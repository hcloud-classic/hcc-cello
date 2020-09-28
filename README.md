[![pipeline status](http://210.207.104.150:8100/iitp-sds/cello/badges/master/pipeline.svg)](http://210.207.104.150:8100/iitp-sds/cello/pipelines)
[![coverage report](http://210.207.104.150:8100/iitp-sds/cello/badges/master/coverage.svg)](http://210.207.104.150:8100/iitp-sds/cello/commits/master)
[![go report](http://210.207.104.150:8100/iitp-sds/hcloud-badge/raw/feature/dev/hcloud-badge_flute.svg)](http://210.207.104.150:8100/iitp-sds/hcloud-badge/raw/feature/dev/goreport_flute)



# Cello 모듈 매뉴얼

```shell
Cello 는 크게 2가지 기능을 담당 합니다.
1. ZFS를 통해 생성한 볼륨을  ISCSI 로 배포
2. pxeboot를 위한 tftp 설정
```





## ~~Graphql~~

```shell
create_volume
update_volume
delete_volume
create_volume_attachment # Additional Volume
update_volume_attachment
delete_volume_attachment
```

현재는 create_volume 만 구현되어있음.



### ~~Exam Query~~

```go
// Volume create
mutation _ {
  create_volume(size:1, filesystem:"zfs", server_uuid:"1234", use_type:"os", user_uuid:"1234") {
    uuid
    size
    filesystem
    server_uuid
    use_type
    user_uuid
    created_at
  }
}
```



## Manage Volume

Storage Node 에서 하는 첫번째 역할은 볼륨 관리이다.

볼륨을 관리는 크게 두가지로 구분된다.

1. 필수 파일시스템 볼륨
2. 데이터 볼륨

두가지의 볼륨으로 구분 하는 이유는 다음과 같다.

### Volume 정의

#### 파일시스템 볼륨

필수 파일시스템 볼륨은 PXEboot 시에 Iscsi로 제공될 파일시스템이 담긴 볼륨을 의미한다. 물론 사전에 관리자에 의해 해당 볼륨은 적절한 운영체제 및 HCC 의 커널및 관련 미들웨어, vnc등 필요한 설정이 전부 완료 되어 있다.

#### 데이터 볼륨

데이터 볼륨은 사용자의 요청에 의해 제공되는 데이터용 볼륨을 의미한다. 쿼터의 제한이 따른다.

데이터 볼륨의 정책은 다음과 같다.

1. 사용자에게 부여된 쿼터의 제한에 따라 생성 요청을 할 수 있다.
2. 생성은 서버 생성 요청 및 추가로 데이터 볼륨을 따로 생성 할 수 있다.
3. 생성된 볼륨의 용량에 대한 수정은 불가 하다.(생성및 삭제만 가능)
4. 볼륨은 파티셔닝 및 포맷이 안되어 있기 때문에 사용자가 직접 해야 한다.

### Volume  배포

생성된 볼륨은 iSCSI 서버를 통해 배포 된다.

각 볼륨은 해당 서버의 UUID를 기준으로 도메인을 설정하여 lun에 등록된다.

#### 주의사항

iSCSI 설정에서 같은 도메인 내에서 볼륨 추가 및 삭제는 `service ctld reload` 를 통해 업데이트를 할 수 있다.

하지만, 볼륨의 크기를 조절했다거나 잘못된 설정이 있을때에는 반드시 `service ctld restart`를 해줘야 한다.

또한 부팅시 특정 볼륨 lun이 인식이 안되고 iSCSI로 login 하는 과정에서 에러가 발생한다면 Storage Node에서 설정이 잘못 되어 있거나, 해당 볼륨이 없기 때문에 발생하는 것이다.

이럴때에는 Storage Node의 설정을 알맞게 수정해 준 후 `service ctld restart`를 해줘야 한다.



## TFTP Server

Storage Node 에서 하는 두번째 역할이다.

### PXESetting

사전에 관리자에 의해 운영체제 별로 pxe 부팅을 위한 기본 설정이 되어 있는 디렉토리가 존재 한다.

해당 디렉토리는 다음 구조를 따르며 최상위 디렉토리는 TFTP에서 설정된 디렉토리이다.

```shell
# /root/boottp 가 최상위 디렉토리 이며 TFTP에 정으되어 있다.
/HCC
|-- defaultCompute
|   |-- initrd.img-2.6.30-hcc-nfs
|   |-- pxelinux.0
|   |-- pxelinux.cfg
|   |   `-- default
|   `-- vmlinuz-2.6.30-hcc
`-- defaultLeader
    |-- initrd.img-2.6.30-hcc
    |-- pxelinux.0
    |-- pxelinux.cfg
    |   `-- default
    `-- vmlinuz-2.6.30-hcc

```

각각의 디렉토리는 Cello에 의해 ServerUUID명의 디렉토리 하위로 복사되며 pxeboot에 필요한 config 파일은 Cello에 의해 내용이 수정된다.

해당 디렉토리(defaultLeader, defaultCompute)는 최초 서버 생성시 서버의 UUID로 구성된 디렉토리에 복사 된다. 설정은 이후에 서버에 맞게 업데이트 된다.

## Clone Volume

볼륨을 Clone 하는 이유는 서버를 생성시 미리 준비된 OS 용 볼륨으로 Provisioning 하기 위함이다.

미리 생성된 OS별 볼륨의 크기는 100G로 고정이며, 이를 기준으로 clone하여 증분데이터만 저장한다. clone 은 원본 볼륨의 스냅샷 단위로 복제가 된다.

## ISCSI Service Handler

추가된 볼륨에 대해 룬 작성

서비스 재시작



config 설정을 하면 reload 해줘야함

### ctl.conf 설정

```shell
$ vim /etc/ctl.conf
portal-group pg0 {
	discovery-auth-group no-authentication
	listen 0.0.0.0
	listen [::]
	
}


lun codex_test {
	path /dev/zvol/master/testimg
	size 1G
}


target iqn.stos.target {
	auth-group no-authentication
	portal-group pg0
	alias codex_target
	lun 0 codex_test
}
```

하지만 cello 에서 다루기 편하게 할라면 시리얼 라이즈 시켜줘야함

```shell

portal-group pg0 { discovery-auth-group no-authentication listen 0.0.0.0 listen [::] }

lun codex_test { path /dev/zvol/master/testimg size 1G }
lun qwe { path /dev/zvol/master/qwe size 1G }
target iqn.stos.target { auth-group no-authentication portal-group pg0 alias codex_target lun 0 codex_test lun 1 qwe }
```

#### 주의사항

ctld 서비스는 **/etc/ctl.conf** 를 파싱할 때, 엄격한 문법을 따진다.

1. 주석 사용 불가

   파일 내용에  주석(comment)을 사용 할 수 업다.

   일반적으로 쉘스크립트에서는 **#** 으로 주석을 쓸 수 있다.

2. 각 요소의 sequential 한 순서

   만약 **iqn.codex.target** 도메인 에 있는 test라는 0번의 **lun** 이 있다고 가정했을 때, 해당 룬 셋팅을 설정값의 순서에서 도메인 보다 밑에 있으면 안된다.

### DB

#### ctl.conf

```
struct iscsitempl{
	target domaintarget []
	
}
struct domaintarget{
	lun string []
}
```



## API List



### CreateVolume

> /RpcCello.Cello/CreateVolume

`Parameter`

```shell
string    UUID = 1; 
string    size = 2; //on demand ex) 100G
string    filesystem = 3; //mandatory ex) centos, ubuntu, debian
string    serverUUID = 4; //mandatory
string    use_type = 5; //mandatory ex)os, data
string    userUUID = 6; //mandatory
google.protobuf.Timestamp created_at = 7;
string    network_IP = 8; //mandatory ex) 172.18.1.1
string    gateway_ip = 9; //mandatory ex) 172.18.1.1
string    pool = 10; //mandatory
int32        lun = 11;
```



#### Example

##### Request

`OS Volume Create`

```shell
size : 100G // 고정 100G를 사용하기 때문에 필요 X
filesystem : centos
serverUUID : xxxx-xxx-xxx-xx-xxx
use_type : os
userUUID : xxx-xxx-xx-xxxx-xx
network_IP : 172.18.1.1 // DHCP IP
gateway_ip : 172.18.1.1 // expose deployment ISCSI IP
pool : volmgmt-1 // User Quota
lun : //생성후 리턴
```

`DATA`

```shell
size : 1024G // 필수
serverUUID : xxxx-xxx-xxx-xx-xxx
use_type : data
userUUID : xxx-xxx-xx-xxxx-xx
network_IP : 172.18.1.1 // DHCP IP
gateway_ip : 172.18.1.1 // expose deployment ISCSI IP
pool : volmgmt-1 // User Quota
lun : //생성후 리턴
```

단, 하단의 상황에 따라 몇가지 제한 사항이 존재.

* 서버 생성시 DATA Volume 추가

  하단 내용 viola추가 개발 해야됨

  ```shell
  서버 생성시 추가된 볼륨은 Viola에 의해 관리자가 지정한 디렉토리로 자동 마운트됨.
  ```

  

* 서버 생성후 DATA Volume 추가



##### Response





## DataSet



`ZFS`

> FileSystem-VolType-ServerUUID



vol_handler

FileSystem : OS(기본 운영체제 볼륨), DATA(추가 볼륨)

OS 볼륨 기본 제공 용량은 100G 고정

```shell
# OS, Data
ex) debian-OS-675acb49-b118-43e9-624a-970139c4f4ff
```



`iscsi`

> TargetName : ServerUUID
>
> Lun : VolType-ServerUUID
>
> 0(OS), 1(DATA)...

object_handler

* iscsi config example

```shell
# iscsi conf /etc/ctl.conf
portal-group pg0 { discovery-auth-group no-authentication listen 0.0.0.0 listen [::] }
lun FileSystem-ServerUUID { path /dev/zvol/volmgmt/FileSystem-VolType-ServerUUID size 100G }
lun FileSystem-ServerUUID { path /dev/zvol/volmgmt/FileSystem-VolType-ServerUUID 5120G }
target iqn.ServerUUID.target {auth-group no-authentication portal-group pg0 lun 0 VolType-ServerUUID lun 1 VolType-ServerUUID}
```







## zfs 관련 command 정리

```
# 이름만 리턴
zpool list -H -o name
master
zpool list -H -o capacity
19%
zpool list -H -o free master
14.1G
zpool list -H -o size master
17.5G
```







### API Example

```
mutation _ {
  create_volume(size:10, filesystem:"centos", server_uuid:"999955fd-08ea-4ede-4e03-f2220724f11d", use_type:"os", user_uuid:"codex",lun_num:0,gateway_ip:"172.18.1.1",pool:"master",network_ip:"172.18.1.1") {
    uuid
    size
    filesystem
    server_uuid
    use_type
    user_uuid
    created_at
  }
  
}

mutation _ {
  create_volume(size:10, filesystem:"centos", server_uuid:"111155fd-08ea-4ede-4e03-f2220724f11d", use_type:"data", user_uuid:"ush",lun_num:0,gateway_ip:"172.18.1.1",pool:"master",network_ip:"172.18.1.1") {
    uuid
    size
    filesystem
    server_uuid
    use_type
    user_uuid
    created_at
  }
  
}
```

