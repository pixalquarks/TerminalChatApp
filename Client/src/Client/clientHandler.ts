import {credentials, ServiceError} from "@grpc/grpc-js";
import {pixalquarks} from "../Proto/chat";
import ServicesClient = pixalquarks.terminalChatServer.ServicesClient;
import ClientName = pixalquarks.terminalChatServer.ClientName;
import FromClient = pixalquarks.terminalChatServer.FromClient;
import {google} from "../Proto/google/protobuf/empty";
import Command = pixalquarks.terminalChatServer.Command;
import Empty = google.protobuf.Empty;
import FromServer = pixalquarks.terminalChatServer.FromServer;
import Clients = pixalquarks.terminalChatServer.Clients;

const sleep = (milliseconds: number) => {
    return new Promise(resolve => setTimeout(resolve, milliseconds))
}


export default class ClientHandler {
    ip : string
    port : string
    name : string = ""
    client
    stream
    constructor(ip:string, port:string) {
        this.ip = ip;
        this.port = port;
        this.client = new ServicesClient(
            `${this.ip}:${this.port}`,
            credentials.createInsecure()
        )
        this.stream = this.client.ChatService()
    }

    ConfigClient(name: string) {
        const cl = new ClientName();
        cl.name = name;
        this.client.VerifyName(cl, (error, resp) => {
            if (error) {
                return error;
            } else {
                this.name = name;
                let msg = new FromClient();
                msg.name = this.name;
                msg.body = "";

                this.stream.write(msg, (err : ServiceError, resp: any) => {
                    if (err) {
                        return err;
                    } else {
                        console.log(resp);
                    }
                })
            }
        });
    }

    OnReceiveMessage(callback: Function) {
        this.stream.on('data', data => callback(data));
    }

    OnEndReceiveMessage(callback: Function) {
        this.stream.on('end', () => callback())
    }

    SendMessage(msg: string, listCallback: Function, PMCallback: Function) {
        console.log(msg);
        const [isCommand, command] = ClientHandler.IsCommand(msg);
        if (isCommand) {
            switch (command) {
                case 'l':
                case 'L':
                    this.ListCommand(listCallback);
                    break;
                case 'p':
                case 'P':
                    this.PMCommand(msg, PMCallback);
                    break;
                default:
                    console.log("not a valid command");
                    break;
            }
        } else {
            let fromClient = new FromClient();
            fromClient.name = this.name;
            fromClient.body = msg;
            this.stream.write(fromClient);
        }
    }

    ListCommand(callback: Function) {
        this.client.GetClients(new google.protobuf.Empty(), (err, data) => callback(err, data));
    }

    PMCommand(msg: string, callback: Function) {
        msg = msg.slice(2);
        let cmd = new Command()
        cmd.type = "P".charCodeAt(0);
        cmd.value = msg;
        cmd.client = this.name;
        this.client.CommandService(cmd, (err: ServiceError|null, value: Empty|undefined) => {
            callback(err, value);
        });
    }

    static IsCommand(msg: string) {
        let isCommand = false;
        if (msg[0] == "!") isCommand = true;
        let cmd = msg[1];
        return [isCommand, cmd];
    }

}


export const LResolve = (err: ServiceError|null, data: Clients|undefined) => {
    if (err) {
        console.log(err);
    } else {
        let clients = data?.client;
        if (clients) {
            for (const c of clients) {
                console.log(c.name);
            }
        }
    }
}

export const PMResolve = (err: ServiceError|null, data: Empty|undefined) => {
    if (err) {
        console.log(err);
    } else {
        console.log(data);
     }
}

export const printMessage = (data:FromServer) => {
    if (data.name == "server") {
        console.log(`******SERVER MESSAGE : ${data.body}******`);
    } else {
        console.log(`${data.name} :: ${data.body}`)
    }
}

export const onStreamEnd = () => {
    console.log("stream ended");
}
