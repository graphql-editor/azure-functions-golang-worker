FROM mcr.microsoft.com/dotnet/core/sdk:3.1 AS runtime-image

ENV PublishWithAspNetCoreTargetManifest=false
ENV HOST_VERSION=3.0.13113
ENV HOST_COMMIT=0cf47580569246787259ef2a29624cf9e8ce61b0

RUN BUILD_NUMBER=$(echo $HOST_VERSION | cut -d'.' -f 3) && \
    wget https://github.com/Azure/azure-functions-host/archive/$HOST_COMMIT.tar.gz && \
    tar xzf $HOST_COMMIT.tar.gz && \
    cd azure-functions-host-* && \
    dotnet publish -v q /p:BuildNumber=$BUILD_NUMBER /p:CommitHash=$HOST_COMMIT src/WebJobs.Script.WebHost/WebJobs.Script.WebHost.csproj --output /azure-functions-host --runtime linux-x64 && \
    rm -rf /azure-functions-host/workers /workers

FROM mcr.microsoft.com/dotnet/core/runtime-deps:3.1 as runtime-deps

COPY --from=runtime-image ["/azure-functions-host", "/azure-functions-host"]

RUN apt-get update && \
    apt-get install -y gnupg wget unzip && \
    wget https://functionscdn.azureedge.net/public/ExtensionBundles/Microsoft.Azure.Functions.ExtensionBundle/1.1.1/Microsoft.Azure.Functions.ExtensionBundle.1.1.1.zip && \
    mkdir -p /FuncExtensionBundles/Microsoft.Azure.Functions.ExtensionBundle/1.1.1 && \
    unzip /Microsoft.Azure.Functions.ExtensionBundle.1.1.1.zip -d /FuncExtensionBundles/Microsoft.Azure.Functions.ExtensionBundle/1.1.1 && \
    rm -f /Microsoft.Azure.Functions.ExtensionBundle.1.1.1.zip

FROM golang:1.13

RUN apt-get update && \
    apt-get install -y libicu63

ENV AzureWebJobsScriptRoot=/home/site/wwwroot \
    HOME=/home \
    FUNCTIONS_WORKER_RUNTIME=golang \
    DOTNET_USE_POLLING_FILE_WATCHER=true \
	ASPNETCORE_URLS=http://*:80

COPY --from=runtime-image [ "/azure-functions-host", "/azure-functions-host" ]
COPY --from=runtime-deps [ "/FuncExtensionBundles", "/FuncExtensionBundles" ]
COPY worker.config.json /azure-functions-host/workers/golang/worker.config.json
EXPOSE 5000
CMD [ "/azure-functions-host/Microsoft.Azure.WebJobs.Script.WebHost" ]
