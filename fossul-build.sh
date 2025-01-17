#!/bin/sh

PLUGIN_DIR="/home/ktenzer/plugins"

if [[ -z "${FOSSUL_BUILD_PLUGIN_DIR}" ]]; then
    export FOSSUL_BUILD_PLUGIN_DIR=$PLUGIN_DIR
fi  

echo "Installing Dependencies"
$GOBIN/dep ensure
#go get github.com/pborman/getopt/v2
#go get github.com/gorilla/mux
#go get k8s.io/client-go/kubernetes/typed/core/v1 
#go get k8s.io/client-go/rest
#go get k8s.io/api/core/v1
#go get k8s.io/client-go/kubernetes/scheme
#go get k8s.io/client-go/tools/remotecommand
#go get github.com/BurntSushi/toml
#go get github.com/lib/pq
#go get github.com/go-sql-driver/mysql
#go get go.mongodb.org/mongo-driver/mongo
#go get go.mongodb.org/mongo-driver/mongo/options
#go get go.mongodb.org/mongo-driver/bson
#go get go.mongodb.org/mongo-driver/mongo/readpref
#go get go.mongodb.org/mongo-driver/x/bsonx

echo "Running Unit Tests"
go test fossul/src/engine/util
if [ $? != 0 ]; then exit 1; fi
go test fossul/src/engine/plugins/pluginUtil
if [ $? != 0 ]; then exit 1; fi

echo "Building Shared Libraries"
go build fossul/src/engine/util
if [ $? != 0 ]; then exit 1; fi
go build fossul/src/engine/client
if [ $? != 0 ]; then exit 1; fi
go build fossul/src/engine/client/k8s
if [ $? != 0 ]; then exit 1; fi
go build fossul/src/engine/plugins/pluginUtil
if [ $? != 0 ]; then exit 1; fi

echo "Building Plugins"
go install fossul/src/engine/plugins/app/basic/sample-app
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $PLUGIN_DIR/app/sample-app.so fossul/src/engine/plugins/app/native/sample-app
if [ $? != 0 ]; then exit 1; fi
go install fossul/src/engine/plugins/storage/basic/sample-storage
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $PLUGIN_DIR/storage/sample-storage.so fossul/src/engine/plugins/storage/native/sample-storage
if [ $? != 0 ]; then exit 1; fi
go install fossul/src/engine/plugins/archive/basic/sample-archive
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $PLUGIN_DIR/archive/sample-archive.so fossul/src/engine/plugins/archive/native/sample-archive
if [ $? != 0 ]; then exit 1; fi
go install fossul/src/engine/plugins/storage/basic/container-basic
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $PLUGIN_DIR/storage/container-basic.so fossul/src/engine/plugins/storage/native/container-basic
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $PLUGIN_DIR/app/mariadb.so fossul/src/engine/plugins/app/native/mariadb
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $PLUGIN_DIR/app/mariadb-dump.so fossul/src/engine/plugins/app/native/mariadb-dump
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $PLUGIN_DIR/app/postgres.so fossul/src/engine/plugins/app/native/postgres
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $PLUGIN_DIR/app/postgres-dump.so fossul/src/engine/plugins/app/native/postgres-dump
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $PLUGIN_DIR/app/mongo.so fossul/src/engine/plugins/app/native/mongo
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $PLUGIN_DIR/app/mongo-dump.so fossul/src/engine/plugins/app/native/mongo-dump
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $PLUGIN_DIR/archive/aws.so fossul/src/engine/plugins/archive/native/aws
if [ $? != 0 ]; then exit 1; fi

echo "Building Services"
go install fossul/src/engine/server
if [ $? != 0 ]; then exit 1; fi
go install fossul/src/engine/app
if [ $? != 0 ]; then exit 1; fi
go install fossul/src/engine/storage
if [ $? != 0 ]; then exit 1; fi

echo "Moving plugins to $PLUGIN_DIR"
if [ ! -d "$PLUGIN_DIR/app" ]; then mkdir $PLUGIN_DIR/app; fi
mv $GOBIN/sample-app $PLUGIN_DIR/app
if [ $? != 0 ]; then exit 1; fi
if [ ! -d "$PLUGIN_DIR/storage" ]; then mkdir $PLUGIN_DIR/storage; fi
mv $GOBIN/sample-storage $PLUGIN_DIR/storage
if [ $? != 0 ]; then exit 1; fi
if [ ! -d "$PLUGIN_DIR/archive" ]; then mkdir $PLUGIN_DIR/archive; fi
mv $GOBIN/sample-archive $PLUGIN_DIR/archive
if [ $? != 0 ]; then exit 1; fi
mv $GOBIN/container-basic $PLUGIN_DIR/storage
if [ $? != 0 ]; then exit 1; fi

echo "Copying default configs"
if [ ! -z $GOBIN/metadata/configs/default ]; then
  mkdir -p $GOBIN/metadata/configs/default
  if [ $? != 0 ]; then exit 1; fi
fi

cp -r $GOPATH/src/fossul/src/cli/configs/default $GOBIN/metadata/configs/default
if [ $? != 0 ]; then exit 1; fi

echo "Copying startup scripts"
cp $GOPATH/src/fossul/fossul-server-startup.sh $GOBIN
if [ $? != 0 ]; then exit 1; fi
cp $GOPATH/src/fossul/fossul-app-startup.sh $GOBIN
if [ $? != 0 ]; then exit 1; fi
cp $GOPATH/src/fossul/fossul-storage-startup.sh $GOBIN
if [ $? != 0 ]; then exit 1; fi

echo "Build completed successfully"

