param location string = 'eastus2'
param name string = 'terminalchatserver'
param containerPort int = 80
param useExternalIngress bool = false
param containerImage string

param envVars array = []

resource law 'Microsoft.OperationalInsights/workspaces@2020-03-01-preview' = {
  name: name
  location: location
  properties: any({
    retentionInDays: 30
    features: {
      searchVersion: 1
    }
    sku: {
      name: 'PerGB2018'
    }
  })
}

resource env 'Microsoft.App/managedEnvironments@2022-03-01' = {
  name: name
  location: location
  properties: {
    appLogsConfiguration: {
      destination: 'log-analytics'
      logAnalyticsConfiguration: {
        customerId: law.properties.customerId
        sharedKey: law.listKeys().primarySharedKey
      }
    }
  }
}

resource containerApp 'Microsoft.App/containerApps@2022-03-01' = {
  name: name
  location: location
  properties: {
    managedEnvironmentId: env.id
    configuration: {
      ingress: {
        external: useExternalIngress
        targetPort: containerPort
        transport: 'http2'
      }
    }
    template: {
      containers: [
        {
          image: containerImage
          name: name
          env: envVars
        }
      ]
    }
  }
}

output fqdn string = containerApp.properties.configuration.ingress.fqdn

// resource containerApp 'Microsoft.Web/containerApps@2021-03-01' = {
//   name: name
//   kind: 'containerapp'
//   location: location
//   properties: {
//     kubeEnvironmentId: containerAppEnvironmentId
//     configuration: {
//       ingress: {
//         external: useExternalIngress
//         targetPort: containerPort
//         transport: 'http2'
//       }
//     }
//     template: {
//       containers: [
//         {
//           image: containerImage
//           name: name
//           env: envVars
//         }
//       ]
//       scale: {
//         minReplicas: 0
//       }
//     }
//   }
// }
