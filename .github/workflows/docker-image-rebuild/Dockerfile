FROM quay.io/ceph/ceph:

RUN dnf -y update && dnf install -y python3-pip xmlsec1 xmlsec1-nss xmlsec1-openssl
RUN pip3 install python3-saml

