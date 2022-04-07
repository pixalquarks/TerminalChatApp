import {pixalquarks} from "../Proto/chat";

const inquirer = require('inquirer');
import ClientHandler from './clientHandler';
import {LResolve, PMResolve, onStreamEnd} from "./clientHandler";
import FromServer = pixalquarks.terminalChatServer.FromServer;

let client: ClientHandler;


inquirer
    .prompt([
        {
            type : 'input',
            name : 'address',
            message : 'Enter the IP Address of the chatroom',
            default : 'localhost:5000',
        },
        {
            type : 'input',
            name : 'port',
            message: 'Enter the port number',
            default: '5000'
        }
    ])
    .then( (answer:any) => {
        client = new ClientHandler(answer.address, answer.port);
        inquirer
            .prompt([
                {
                    type : 'input',
                    name : 'name',
                    message : 'Please enter your name',
                }
            ])
            .then( (answer:any) => {
                client.ConfigClient(answer.name);
                client.OnReceiveMessage(printMessage);
                client.OnEndReceiveMessage(onStreamEnd);
                chat()
                    .then((val) => {
                        console.log('chat ended');
                    })
                    .catch((err) => {
                        console.log(err);
                    })
            })
    })
    .catch ((error:any) => {
        console.log(error)
    })

async function chat() {
    while (true) {
        const msg = await inquirer.prompt({
            name: "sendMsg",
            type: "input",
            message: "-->",
        });
        if (msg.sendMsg == '!q') break;
        if (msg.sendMsg === '!c') console.clear();
        client.SendMessage(msg.sendMsg, LResolve, PMResolve);
    }
}

const printMessage = (data:FromServer) => {
    if (data.name == "server") {
        console.log(`******SERVER MESSAGE : ${data.body}******`);
    } else {
        console.log(`${data.name} :: ${data.body}`);
    }
}
