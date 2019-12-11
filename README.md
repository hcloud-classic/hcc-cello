[![pipeline status](http://210.207.104.140:8100/iitp-sds/cello/badges/feature/toflute/pipeline.svg)](http://210.207.104.140:8100/iitp-sds/cello/pipelines)
[![coverage report](http://210.207.104.140:8100/iitp-sds/cello/badges/feature/toflute/coverage.svg)](http://210.207.104.140:8100/iitp-sds/cello/commits/feature/toflute)
[![go report](http://210.207.104.140:8100/iitp-sds/hcloud-badge/raw/feature/dev/hcloud-badge_flute.svg)](http://210.207.104.140:8100/iitp-sds/hcloud-badge/raw/feature/dev/goreport_flute)





# Cello 모듈 매뉴얼



## Graphql

```shell
create_volume
update_volume
delete_volume
create_volume_attachment # Additional Volume
update_volume_attachment
delete_volume_attachment
```

현재는 create_volume 만 구현되어있음.



### Exam Query

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



## Create Volume

볼륨을 생성은 크게 두가지로 구분된다.

1. 필수 파일시스템 볼륨
2. 데이터 볼륨

두가지의 볼륨으로 구분지은 이유는 다음과 같다.

### 볼륨 타입

#### 파일시스템 볼륨

필수 파일시스템 볼륨은 PXEboot 시에 Iscsi로 제공될 파일시스템이 담긴 볼륨을 의미한다. 물론 사전에 관리자에 의해 해당 볼륨은 적절한 운영체제 및 HCC 의 커널로 설치가 완료 되어 있다.

#### 데이터 볼륨

데이터 볼륨은 사용자의 요청에 의해 제공되는 데이터용 볼륨을 의미한다. 쿼터의 제한이 따른다.

볼륨 타입을 정의 했다면 다음은 PXEBoot 를 하기위한 설정과 데이터의 준비이다.

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



