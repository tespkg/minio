module github.com/minio/minio

go 1.19

require (
	cloud.google.com/go/storage v1.28.1
	github.com/Azure/azure-sdk-for-go v27.0.0+incompatible
	github.com/Azure/go-autorest v11.7.0+incompatible
	github.com/alecthomas/participle v0.2.1
	github.com/aliyun/aliyun-oss-go-sdk v0.0.0-20190307165228-86c17b95fcd5
	github.com/bcicen/jstream v0.0.0-20190220045926-16c1f8af81c2
	github.com/cheggaaa/pb v1.0.28
	github.com/cncf/udpa/go v0.0.0-20220112060539-c52dc94e7fbe
	github.com/coredns/coredns v1.4.0
	github.com/coreos/go-semver v0.2.0
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/djherbis/atime v1.0.0
	github.com/dustin/go-humanize v1.0.0
	github.com/eclipse/paho.mqtt.golang v1.1.2-0.20190322152051-20337d8c3947
	github.com/elazarl/go-bindata-assetfs v1.0.0
	github.com/envoyproxy/go-control-plane v0.11.1-0.20230406144219-ba92d50b6596
	github.com/fatih/color v1.7.0
	github.com/fatih/structs v1.1.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/gogo/protobuf v1.3.2
	github.com/golang/glog v1.0.0
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e
	github.com/golang/protobuf v1.5.3
	github.com/golang/snappy v0.0.1
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/btree v0.0.0-20180813153112-4030bb1f1f0c
	github.com/google/go-cmp v0.5.9
	github.com/gorilla/handlers v1.4.0
	github.com/gorilla/mux v1.7.0
	github.com/gorilla/rpc v1.2.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-gateway v1.7.0
	github.com/hashicorp/vault v1.1.0
	github.com/inconshreveable/go-update v0.0.0-20160112193335-8152e7eb6ccf
	github.com/jonboulle/clockwork v0.1.0
	github.com/json-iterator/go v1.1.9
	github.com/klauspost/pgzip v1.2.1
	github.com/klauspost/readahead v1.3.0
	github.com/klauspost/reedsolomon v1.9.1
	github.com/kr/pty v1.1.3
	github.com/lib/pq v1.0.0
	github.com/mattn/go-isatty v0.0.7
	github.com/miekg/dns v1.1.8
	github.com/minio/blazer v0.0.0-20171126203752-2081f5bf0465
	github.com/minio/cli v1.3.0
	github.com/minio/dsync v1.0.0
	github.com/minio/hdfs/v3 v3.0.0
	github.com/minio/highwayhash v1.0.2
	github.com/minio/lsync v0.0.0-20190207022115-a4e43e3d0887
	github.com/minio/mc v0.0.0-20190401030144-a1355e50e2e8
	github.com/minio/minio-go/v6 v6.0.58-0.20200611233712-3da267908ecf
	github.com/minio/parquet-go v0.0.0-20190318185229-9d767baf1679
	github.com/minio/sha256-simd v0.1.1
	github.com/minio/sio v0.0.0-20190118043801-035b4ef8c449
	github.com/mitchellh/go-homedir v1.1.0
	github.com/nats-io/nats.go v1.26.0
	github.com/nats-io/stan.go v0.10.4
	github.com/ncw/directio v1.0.5
	github.com/nsqio/go-nsq v1.0.7
	github.com/pkg/profile v1.3.0
	github.com/prometheus/client_golang v0.9.3-0.20190127221311-3c4408c8b829
	github.com/prometheus/client_model v0.3.0
	github.com/rjeczalik/notify v0.9.2
	github.com/rs/cors v1.6.0
	github.com/segmentio/go-prompt v1.2.1-0.20161017233205-f0d19b6901ad
	github.com/skyrings/skyring-common v0.0.0-20160929130248-d1c0bb1cbd5e
	github.com/soheilhy/cmux v0.1.4
	github.com/spf13/pflag v1.0.3
	github.com/streadway/amqp v0.0.0-20190402114354-16ed540749f6
	github.com/tidwall/gjson v1.2.1
	github.com/tidwall/sjson v1.0.4
	github.com/tmc/grpc-websocket-proxy v0.0.0-20190109142713-0ad062ec5ee5
	github.com/ugorji/go/codec v0.0.0-20190320090025-2dc34c0b8780
	github.com/valyala/tcplisten v0.0.0-20161114210144-ceec8f93295a
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2
	go.uber.org/atomic v1.3.2
	go.uber.org/zap v1.9.1
	golang.org/x/crypto v0.8.0
	golang.org/x/net v0.9.0
	golang.org/x/oauth2 v0.7.0
	golang.org/x/sys v0.7.0
	golang.org/x/time v0.3.0
	google.golang.org/api v0.114.0
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1
	gopkg.in/Shopify/sarama.v1 v1.20.0
	gopkg.in/olivere/elastic.v5 v5.0.80
	gopkg.in/yaml.v2 v2.2.2
)

require (
	cloud.google.com/go v0.110.0 // indirect
	cloud.google.com/go/compute v1.19.1 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/iam v0.13.0 // indirect
	contrib.go.opencensus.io/exporter/ocagent v0.4.10 // indirect
	git.apache.org/thrift.git v0.12.0 // indirect
	github.com/DataDog/zstd v1.3.5 // indirect
	github.com/armon/go-metrics v0.0.0-20190430140413-ec5e00d3c878 // indirect
	github.com/beorn7/perks v0.0.0-20180321164747-3a771d992973 // indirect
	github.com/census-instrumentation/opencensus-proto v0.4.1 // indirect
	github.com/cncf/xds/go v0.0.0-20230310173818-32f1caf87195 // indirect
	github.com/coreos/go-systemd v0.0.0-20190321100706-95778dfbb74e // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/eapache/go-resiliency v1.1.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20180814174437-776d5712da21 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/envoyproxy/protoc-gen-validate v0.10.1 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.3 // indirect
	github.com/googleapis/gax-go/v2 v2.7.1 // indirect
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.11.3 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.0 // indirect
	github.com/hashicorp/go-multierror v1.0.0 // indirect
	github.com/hashicorp/go-retryablehttp v0.5.3 // indirect
	github.com/hashicorp/go-rootcerts v0.0.0-20160503143440-6bb64b370b90 // indirect
	github.com/hashicorp/go-sockaddr v1.0.2 // indirect
	github.com/hashicorp/go-uuid v1.0.1 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/howeyc/gopass v0.0.0-20170109162249-bf9dde6d0d2c // indirect
	github.com/jcmturner/gofork v0.0.0-20180107083740-2aebee971930 // indirect
	github.com/klauspost/compress v1.16.5 // indirect
	github.com/klauspost/cpuid v1.2.3 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.1 // indirect
	github.com/mailru/easyjson v0.0.0-20180730094502-03f2033d19d5 // indirect
	github.com/marstr/guid v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.1 // indirect
	github.com/mattn/go-runewidth v0.0.4 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/minio/md5-simd v1.1.0 // indirect
	github.com/minio/minio-go v0.0.0-20190313212832-5d20267d970d // indirect
	github.com/mitchellh/mapstructure v0.0.0-20160808181253-ca63d7c062ee // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/nats-io/go-nats v1.7.2 // indirect
	github.com/nats-io/go-nats-streaming v0.4.2 // indirect
	github.com/nats-io/nats-server/v2 v2.9.17 // indirect
	github.com/nats-io/nkeys v0.4.4 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pierrec/lz4 v2.0.5+incompatible // indirect
	github.com/pkg/errors v0.8.1 // indirect
	github.com/prometheus/common v0.2.0 // indirect
	github.com/prometheus/procfs v0.0.0-20190117184657-bf6a532e95b1 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20181016184325-3113b8401b8a // indirect
	github.com/ryanuber/go-glob v1.0.0 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/sirupsen/logrus v1.5.0 // indirect
	github.com/tidwall/match v1.0.1 // indirect
	github.com/tidwall/pretty v0.0.0-20190325153808-1166b9ac2b65 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/term v0.7.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/grpc v1.54.0 // indirect
	google.golang.org/grpc/examples v0.0.0-20230524173754-59134c303c31 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/ini.v1 v1.42.0 // indirect
	gopkg.in/jcmturner/aescts.v1 v1.0.1 // indirect
	gopkg.in/jcmturner/dnsutils.v1 v1.0.1 // indirect
	gopkg.in/jcmturner/goidentity.v2 v2.0.0 // indirect
	gopkg.in/jcmturner/gokrb5.v5 v5.3.0 // indirect
	gopkg.in/jcmturner/rpc.v0 v0.0.2 // indirect
)
