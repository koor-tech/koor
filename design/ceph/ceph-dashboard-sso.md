# Adding Ceph Dashboard SSO support

## Goal
- Perform Ceph Dashboard SSO setup using KSD/Rook

## Prerequisites
- Ceph Dashboard enabled
- Enable SSO and IDP client creation (e.g., in Keycloak it is called a client, other software might call it differently, but that is out of scope for this feature)

## Requirements
- [koor-ceph-container](https://hub.docker.com/layers/koorinc/koor-ceph-container) image from docker hub.
- Ceph command to enable the dashboard: `# ceph dashboard sso setup saml2 <ceph_dashboard_base_url> <idp_metadata> {<idp_username_attribute>} {<idp_entity_id>} {<sp_x_509_cert>} {<sp_private_key>}`

## What we need to accomplish
The following parameters need to be added in the CephCluster CRD to enable SSO:
- sso: # By default it’ll be commented and disabled.
  - enabled: true|false
  - If enabled then we need the following parameters:
    - ceph_dashboard_base_url
    - Idp_metadata
  - The following parameters are optional requirements:
    - `<idp_username_attribute>` # Default value is “username” if I saw it in the code correctly, so technically optional value.
    - `<idp_entity_id>` # We can auto generate this based on the `ceph_dashboard_base_url`, so it can be made optional.
    - `<sp_x_509_cert>`
    - `<sp_private_key>` # We would require both the cert and private key to be set when one of them is set.

## Result

```yaml
apiVersion: ceph.rook.io/v1
kind: CephCluster
metadata:
  name: rook-ceph
  namespace: rook-ceph # namespace:cluster
spec:
[...]
    dashboard:
    enabled: true
    # serve the dashboard under a subpath (useful when you are accessing the dashboard via a reverse proxy)
    # urlPrefix: /ceph-dashboard
    # serve the dashboard at the given port.
    # port: 8443
    # serve the dashboard using SSL
    ssl: true
    # configure sso for the dashboard access
    sso:
      enabled: false # by default it’ll be false
      users:
        - username:
          roles:
           - <role>
      baseUrl: <url_value>
      idpMetadataUrl: <url_value>
      idpUsernameAttribute: "username"
      entityID: <entity_id_value>
      spCert:
        key: tls.crt
        secret:
        name: my-secret
      spPrivateKey:
        key: tls.key
        secret:
          name: my-secret
```

Use https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#secretreference-v1-core or something similar to target a secret in the same namespace for the spCert.secret and spPrivateKey.secret

**Note:**

The cert and key need to be mounted into the MGR instances.

The issuer value of SAML requests will follow this pattern: `<ceph_dashboard_base_url>/auth/saml2/metadata`

## References

* https://docs.ceph.com/en/latest/mgr/dashboard/#enabling-single-sign-on-sso
* https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/
