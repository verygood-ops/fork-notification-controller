# Providers

<!-- menuweight:40 -->

The `Provider` API defines how events are encoded and where to send them.

## Example

The following is an example of how to send alerts to Slack when Flux fails to
install or upgrade [Flagger](https://github.com/fluxcd/flagger).

```yaml
---
apiVersion: notification.toolkit.fluxcd.io/v1beta3
kind: Provider
metadata:
  name: slack-bot
  namespace: flagger-system
spec:
  type: slack
  channel: general
  address: https://slack.com/api/chat.postMessage
  secretRef:
    name: slack-bot-token
---
apiVersion: notification.toolkit.fluxcd.io/v1beta3
kind: Alert
metadata:
  name: slack
  namespace: flagger-system
spec:
  summary: "Flagger impacted in us-east-2"
  providerRef:
    name: slack-bot
  eventSeverity: error
  eventSources:
    - kind: HelmRepository
      name: '*'
    - kind: HelmRelease
      name: '*'
```

In the above example:

- A Provider named `slack-bot` is created, indicated by the
  `Provider.metadata.name` field.
- An Alert named `slack` is created, indicated by the
  `Alert.metadata.name` field.
- The Alert references the `slack-bot` provider, indicated by the
  `Alert.spec.providerRef` field.
- The notification-controller starts listening for events sent for
  all HelmRepositories and HelmReleases in the `flagger-system` namespace.
- When an event with severity `error` is received, the controller posts
  a message on Slack containing the `summary` text and the Helm install or
  upgrade error.
- The controller uses the Slack Bot token from the secret indicated by the
  `Provider.spec.secretRef.name` to authenticate with the Slack API.

You can run this example by saving the manifests into `slack-alerts.yaml`.

1. First create a secret with the Slack bot token:

   ```sh
   kubectl -n flagger-system create secret generic slack-bot-token --from-literal=token=xoxb-YOUR-TOKEN
   ```

2. Apply the resources on the cluster:

   ```sh
   kubectl -n flagger-system apply --server-side -f slack-alerts.yaml
   ```

## Writing a provider spec

As with all other Kubernetes config, a Provider needs `apiVersion`,
`kind`, and `metadata` fields. The name of an Alert object must be a
valid [DNS subdomain name](https://kubernetes.io/docs/concepts/overview/working-with-objects/names#dns-subdomain-names).

A Provider also needs a
[`.spec` section](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#spec-and-status).

### Type

`.spec.type` is a required field that specifies which SaaS API to use.

The supported alerting providers are:

| Provider                                                | Type             |
|---------------------------------------------------------|------------------|
| [Generic webhook](#generic-webhook)                     | `generic`        |
| [Generic webhook with HMAC](#generic-webhook-with-hmac) | `generic-hmac`   |
| [Azure Event Hub](#azure-event-hub)                     | `azureeventhub`  |
| [DataDog](#datadog)                                     | `datadog`        |
| [Discord](#discord)                                     | `discord`        |
| [GitHub dispatch](#github-dispatch)                     | `githubdispatch` |
| [Google Chat](#google-chat)                             | `googlechat`     |
| [Google Pub/Sub](#google-pubsub)                        | `googlepubsub`   |
| [Grafana](#grafana)                                     | `grafana`        |
| [Lark](#lark)                                           | `lark`           |
| [Matrix](#matrix)                                       | `matrix`         |
| [Microsoft Teams](#microsoft-teams)                     | `msteams`        |
| [Opsgenie](#opsgenie)                                   | `opsgenie`       |
| [PagerDuty](#pagerduty)                                 | `pagerduty`      |
| [Prometheus Alertmanager](#prometheus-alertmanager)     | `alertmanager`   |
| [Rocket](#rocket)                                       | `rocket`         |
| [Sentry](#sentry)                                       | `sentry`         |
| [Slack](#slack)                                         | `slack`          |
| [Telegram](#telegram)                                   | `telegram`       |
| [WebEx](#webex)                                         | `webex`          |
| [NATS](#nats)                                           | `nats`           |

#### Types supporting Git commit status updates

The providers supporting [Git commit status updates](#git-commit-status-updates) are:

| Provider                                                     | Type              |
|--------------------------------------------------------------|-------------------|
| [Azure DevOps](#azure-devops)                                | `azuredevops`     |
| [Bitbucket](#bitbucket)                                      | `bitbucket`       |
| [Bitbucket Server/Data Center](#bitbucket-serverdata-center) | `bitbucketserver` |
| [GitHub](#github)                                            | `github`          |
| [GitLab](#gitlab)                                            | `gitlab`          |
| [Gitea](#gitea)                                              | `gitea`           |

#### Alerting

##### Generic webhook

When `.spec.type` is set to `generic`, the controller will send an HTTP POST
request to the provided [Address](#address).

The body of the request is a [JSON `Event` object](events.md#event-structure),
for example:

```json
{
  "involvedObject": {
    "apiVersion": "kustomize.toolkit.fluxcd.io/v1",
    "kind": "Kustomization",
    "name": "webapp",
    "namespace": "apps",
    "uid": "7d0cdc51-ddcf-4743-b223-83ca5c699632"
  },
  "metadata": {
    "kustomize.toolkit.fluxcd.io/revision": "main/731f7eaddfb6af01cb2173e18f0f75b0ba780ef1"
  },
  "severity":"error",
  "reason": "ValidationFailed",
  "message":"service/apps/webapp validation error: spec.type: Unsupported value: Ingress",
  "reportingController":"kustomize-controller",
  "reportingInstance":"kustomize-controller-7c7b47f5f-8bhrp",
  "timestamp":"2022-10-28T07:26:19Z"
}
```

Where the `involvedObject` key contains the metadata from the object triggering
the event.

The controller includes a `Gotk-Component` header in the request, which can be
used to identify the component which sent the event, e.g. `source-controller`
or `notification-controller`.

```
POST / HTTP/1.1
Host: example.com
Accept-Encoding: gzip
Content-Length: 452
Content-Type: application/json
Gotk-Component: kustomize-controller
User-Agent: Go-http-client/1.1
```

You can add additional headers to the POST request using a [`headers` key in the
referenced Secret](#http-headers-example).

##### Generic webhook with HMAC

When `.spec.type` is set to `generic-hmac`, the controller will send an HTTP
POST request to the provided [Address](#address) for an [Event](events.md#event-structure),
while including an `X-Signature` HTTP header carrying the HMAC of the request
body. The inclusion of the header allows the receiver to verify the
authenticity and integrity of the request.

The `X-Signature` header is calculated by generating an HMAC of the request
body using the [`token` key from the referenced Secret](#token-example). The
HTTP header value has the following format:

```
X-Signature: <hash-function>=<hash>
```

`<hash-function>` denotes the hash function used to generate the HMAC and
currently defaults to `sha256`, which may change in the future. `<hash>` is the
HMAC of the request body, encoded as a hexadecimal string.

while `<hash>` is the hex-encoded HMAC value.

The body of the request is a [JSON `Event` object](events.md#event-structure),
as described in the [Generic webhook](#generic-webhook) section.

###### HMAC verification example

The following example in Go shows how to verify the authenticity and integrity
of a request by using the X-Signature header.

```go
func verifySignature(signature string, payload, key []byte) error {
	sig := strings.Split(signature, "=")

	if len(sig) != 2 {
		return fmt.Errorf("invalid signature value")
	}

	var newF func() hash.Hash
	switch sig[0] {
	case "sha224":
		newF = sha256.New224
	case "sha256":
		newF = sha256.New
	case "sha384":
		newF = sha512.New384
	case "sha512":
		newF = sha512.New
	default:
		return fmt.Errorf("unsupported signature algorithm %q", sig[0])
	}

	mac := hmac.New(newF, key)
	if _, err := mac.Write(payload); err != nil {
		return fmt.Errorf("failed to write payload to HMAC encoder: %w", err)
	}

	sum := fmt.Sprintf("%x", mac.Sum(nil))
	if sum != sig[1] {
		return fmt.Errorf("HMACs do not match: %#v != %#v", sum, sig[1])
	}
	return nil
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Require a X-Signature header
	if len(r.Header["X-Signature"]) == 0 {
		http.Error(w, "missing X-Signature header", http.StatusBadRequest)
		return
	}

	// Read the request body with a limit of 1MB
	lr := io.LimitReader(r.Body, 1<<20)
	body, err := io.ReadAll(lr)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	// Verify signature using the same token as the Secret referenced in
	// Provider
	key := []byte("<token>")
	if err := verifySignature(r.Header.Get("X-Signature"), body, key); err != nil {
		http.Error(w, fmt.Sprintf("failed to verify HMAC signature: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// Do something with the verified request body
	// ...
}
```

##### Slack

When `.spec.type` is set to `slack`, the controller will send a message for an
[Event](events.md#event-structure) to the provided Slack API [Address](#address). 

The Event will be formatted into a Slack message using an [Attachment](https://api.slack.com/reference/messaging/attachments),
with the metadata attached as fields, and the involved object as author.
The severity of the Event is used to set the color of the attachment.

When a [Channel](#channel) is provided, it will be added as a [`channel`
field](https://api.slack.com/methods/chat.postMessage#arg_channel) to the API
payload. Otherwise, the further configuration of the [Address](#address) will
determine the channel.

When [Username](#username) is set, this will be added as a [`username`
field](https://api.slack.com/methods/chat.postMessage#arg_username) to the
payload, defaulting to the name of the reporting controller.

This Provider type supports the configuration of a [proxy URL](#https-proxy)
and/or [certificate secret reference](#certificate-secret-reference).

###### Slack example

To configure a Provider for Slack, we recommend using a Slack Bot App token which is
not attached to a specific Slack account. To obtain a token, please follow
[Slack's guide on creating an app](https://api.slack.com/authentication/basics#creating).

Once you have obtained a token, [create a Secret containing the `token`
key](#token-example) and a `slack` Provider with the `address` set to
`https://slack.com/api/chat.postMessage`.

Using this API endpoint, it is possible to send messages to multiple channels
by adding the integration to each channel.

```yaml
---
apiVersion: notification.toolkit.fluxcd.io/v1beta3
kind: Provider
metadata:
  name: slack
  namespace: default
spec:
  type: slack
  channel: general
  address: https://slack.com/api/chat.postMessage
  secretRef:
    name: slack-token
---
apiVersion: v1
kind: Secret
metadata:
  name: slack-token
  namespace: default
stringData:
    token: xoxb-1234567890-1234567890-1234567890-1234567890
```

###### Slack (legacy) example

To configure a Provider for Slack using the [legacy incoming webhook API](https://api.slack.com/messaging/webhooks),
create a Secret with the `address` set to `https://hooks.slack.com/services/...`,
and a `slack` Provider with a [Secret reference](#address-example).

```yaml
---
apiVersion: notification.toolkit.fluxcd.io/v1beta3
kind: Provider
metadata:
  name: slack
  namespace: default
spec:
  type: slack
  secretRef:
    name: slack-webhook
---
apiVersion: v1
kind: Secret
metadata:
  name: slack-webhook
  namespace: default
stringData:
  address: "https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK"
```

##### Microsoft Teams

When `.spec.type` is set to `msteams`, the controller will send a payload for
an [Event](events.md#event-structure) to the provided [Address](#address). The address
may be a [Microsoft Teams Incoming Webhook Workflow](https://support.microsoft.com/en-us/office/create-incoming-webhooks-with-workflows-for-microsoft-teams-8ae491c7-0394-4861-ba59-055e33f75498), or
the deprecated [Office 365 Connector](https://devblogs.microsoft.com/microsoft365dev/retirement-of-office-365-connectors-within-microsoft-teams/).

**Note:** If the Address host contains the suffix `.webhook.office.com`, the controller will imply that
the backend is the deprecated Office 365 Connector and is expecting the Event in the [connector message](https://learn.microsoft.com/en-us/microsoftteams/platform/webhooks-and-connectors/how-to/connectors-using#example-of-connector-message) format. Otherwise, the controller will format the Event as a [Microsoft Adaptive Card](https://adaptivecards.io/explorer/) message.

In both cases the Event metadata is attached as facts, and the involved object as a summary/title.
The severity of the Event is used to set the color of the message.

This Provider type supports the configuration of a [proxy URL](#https-proxy)
and/or [certificate secret reference](#certificate-secret-reference), but lacks support for
configuring a [Channel](#channel). This can be configured during the
creation of the Incoming Webhook Workflow in Microsoft Teams.

###### Microsoft Teams example

To configure a Provider for Microsoft Teams, create a Secret with [the
`address`](#address-example) set to the [webhook URL](https://support.microsoft.com/en-us/office/create-incoming-webhooks-with-workflows-for-microsoft-teams-8ae491c7-0394-4861-ba59-055e33f75498),
and an `msteams` Provider with a [Secret reference](#secret-reference).

```yaml
---
apiVersion: notification.toolkit.fluxcd.io/v1beta3
kind: Provider
metadata:
  name: msteams
  namespace: default
spec:
  type: msteams
  secretRef:
    name: msteams-webhook
---
apiVersion: v1
kind: Secret
metadata:
  name: msteams-webhook
  namespace: default
stringData:
  address: https://prod-xxx.yyy.logic.azure.com:443/workflows/zzz/triggers/manual/paths/invoke?...
```

##### DataDog

When `.spec.type` is set to `datadog`, the controller will send a payload for
an [Event](events.md#event-structure) to the provided DataDog API [Address](#address).

The Event will be formatted into a [DataDog Event](https://docs.datadoghq.com/api/latest/events/#post-an-event) and sent to the
API endpoint of the provided DataDog [Address](#address).

This Provider type supports the configuration of a [proxy URL](#https-proxy)
and/or [certificate secret reference](#certificate-secret-reference).

The metadata of the Event is included in the DataDog event as extra tags.

###### DataDog example

To configure a Provider for DataDog, create a Secret with [the `token`](#token-example)
set to a [DataDog API key](https://docs.datadoghq.com/account_management/api-app-keys/#api-keys)
(not an application key!) and a `datadog` Provider with a [Secret reference](#secret-reference).

```yaml
---
apiVersion: notification.toolkit.fluxcd.io/v1beta3
kind: Provider
metadata:
  name: datadog
  namespace: default
spec:
  type: datadog
  address: https://api.datadoghq.com # DataDog Site US1
  secretRef:
    name: datadog-secret
---
apiVersion: v1
kind: Secret
metadata:
  name: datadog-secret
  namespace: default
stringData:
  token: <DataDog API Key>
---
apiVersion: notification.toolkit.fluxcd.io/v1beta3
kind: Alert
metadata:
  name: datadog-info
  namespace: default
spec:
  eventSeverity: info
  eventSources:
    - kind: HelmRelease
      name: "*"
  providerRef:
    name: datadog
  eventMetadata:
    env: my-k8s-cluster # example of adding a custom `env` tag to the event
```

##### Discord

When `.spec.type` is set to `discord`, the controller will send a payload for
an [Event](events.md#event-structure) to the provided Discord [Address](#address).

The Event will be formatted into a [Slack message](#slack) and send to the
`/slack` endpoint of the provided Discord [Address](#address).

This Provider type supports the configuration of a [proxy URL](#https-proxy)
and/or [certificate secret reference](#certificate-secret-reference), but lacks support for
configuring a [Channel](#channel). This can be configured [during the creation
of the address](https://discord.com/developers/docs/resources/webhook#create-webhook)

###### Discord example

To configure a Provider for Discord, create a Secret with [the `address`](#address-example)
set to the [webhook URL](https://discord.com/developers/docs/resources/webhook#create-webhook),
and a `discord` Provider with a [Secret reference](#secret-reference).

```yaml
---
apiVersion: notification.toolkit.fluxcd.io/v1beta3
kind: Provider
metadata:
  name: discord
  namespace: default
spec:
  type: discord
  secretRef:
    name: discord-webhook
---
apiVersion: v1
kind: Secret
metadata:
  name: discord-webhook
  namespace: default
stringData:
    address: "https://discord.com/api/webhooks/..."
```


##### Sentry

When `.spec.type` is set to `sentry`, the controller will send a payload for
an [Event](events.md#event-structure) to the provided Sentry [Address](#address).

Depending on the `severity` of the Event, the controller will capture a [Sentry
Event](https://develop.sentry.dev/sdk/event-payloads/)for `error`, or [Sentry
Transaction Event](https://develop.sentry.dev/sdk/event-payloads/transaction/)
with a [Span](https://develop.sentry.dev/sdk/event-payloads/span/) for `info`.
The metadata of the Event is included as [`extra` data](https://develop.sentry.dev/sdk/event-payloads/#optional-attributes)
in the Sentry Event, or as [Span `tags`](https://develop.sentry.dev/sdk/event-payloads/span/#attributes).

The Provider's [Channel](#channel) is used to set the `environment` on the
Sentry client.

This Provider type supports the configuration of
[certificate secret reference](#certificate-secret-reference).

###### Sentry example

To configure a Provider for Sentry, create a Secret with [the `address`](#address-example)
set to a [Sentry DSN](https://docs.sentry.io/product/sentry-basics/dsn-explainer/),
and a `sentry` Provider with a [Secret reference](#secret-reference).

```yaml
---
apiVersion: notification.toolkit.fluxcd.io/v1beta3
kind: Provider
metadata:
  name: sentry
  namespace: default
spec:
  type: sentry
  channel: staging-env
  secretRef:
    name: sentry-webhook
---
apiVersion: v1
kind: Secret
metadata:
  name: sentry-webhook
  namespace: default
stringData:
  address: "https://....@sentry.io/12341234"
```

**Note:** The `sentry` Provider also sends traces for events with the severity
`info`. This can be disabled by setting the `spec.eventSeverity` field to
`error` on an `Alert`.

##### Telegram

When `.spec.type` is set to `telegram`, the controller will send a payload for
an [Event](events.md#event-structure) to the Telegram Bot API.

The Event will be formatted into a message string, with the metadata attached
as a list of key-value pairs.

The Provider's [Channel](#channel) is used to set the receiver of the message.
This can be a unique identifier (`-1234567890`) for the target chat,
a unique identifier with the topic identifier (`-1234567890:1`) for the forum chat,
or the username (`@username`) of the target channel.

This Provider type does not support the configuration of a [certificate secret reference](#certificate-secret-reference).

**Note:** The Telegram provider always ignores the `address` field and uses the default 
Telegram Bot API endpoint (`https://api.telegram.org`).

###### Telegram example

To configure a Provider for Telegram, create a Secret with [the `token`](#token-example)
obtained from [the BotFather](https://core.telegram.org/bots#how-do-i-create-a-bot),
and a `telegram` Provider with a [Secret reference](#secret-reference).

```yaml
---
apiVersion: notification.toolkit.fluxcd.io/v1beta3
kind: Provider
metadata:
  name: telegram
  namespace: default
spec:
  type: telegram
  channel: "@fluxcd" # or "-1557265138" (channel id) or "-1552289257:1" (forum chat id with topic id)
  secretRef:
    name: telegram-token
```

##### Matrix

When `.spec.type` is set to `matrix`, the controller will send a payload for
an [Event](events.md#event-structure) to the provided Matrix [Address](#address).

The Event will be formatted into a message string, with the metadata attached 
as a list of key-value pairs, and send as a [`m.room.message` text event](https://spec.matrix.org/v1.3/client-server-api/#mroommessage)
to the provided Matrix [Address](#address).

The Provider's [Channel](#channel) is used to set the receiver of the message
using a room identifier (`!1234567890:example.org`).

This provider type does support the configuration of [TLS
certificates](#tls-certificates).

###### Matrix example

To configure a Provider for Matrix, create a Secret with [the `token`](#token-example)
obtained from [the Matrix endpoint](https://matrix.org/docs/guides/client-server-api#registration),
and a `matrix` Provider with a [Secret reference](#secret-reference).

```yaml
apiVersion: notification.toolkit.fluxcd.io/v1beta3
kind: Provider
metadata:
  name: matrix
  namespace: default
spec:
  type: matrix
  address: https://matrix.org
  channel: "!jezptmDwEeLapMLjOc:matrix.org"
  secretRef:
    name: matrix-token
```

##### Lark

When `.spec.type` is set to `lark`, the controller will send a payload for
an [Event](events.md#event-structure) to the provided Lark [Address](#address).

The Event will be formatted into a [Lark Message card](https://open.larksuite.com/document/ukTMukTMukTM/uczM3QjL3MzN04yNzcDN),
with the metadata written to the message string.

This Provider type does not support the configuration of a [proxy URL](#https-proxy)
or [certificate secret reference](#certificate-secret-reference).

###### Lark example

To configure a Provider for Lark, create a Secret with [the `address`](#address-example)
obtained from [adding a bot to a group](https://open.larksuite.com/document/uAjLw4CM/ukTMukTMukTM/bot-v3/use-custom-bots-in-a-group#57181e84),
and a `lark` Provider with a [Secret reference](#secret-reference).

```yaml
apiVersion: notification.toolkit.fluxcd.io/v1beta3
kind: Provider
metadata:
  name: lark
  namespace: default
spec:
  type: lark
  secretRef:
    name: lark-webhook
---
apiVersion: v1
kind: Secret
metadata:
  name: lark-webhook
  namespace: default
stringData:
    address: "https://open.larksuite.com/open-apis/bot/v2/hook/xxxxxxxxxxxxxxxxx"
```

##### Rocket

When `.spec.type` is set to `rocket`, the controller will send a payload for
an [Event](events.md#event-structure) to the provided Rocket [Address](#address).

The Event will be formatted into a [Slack message](#slack) and send as a
payload the provided Rocket [Address](#address).

This Provider type does support the configuration of a [proxy URL](#https-proxy)
and [certificate secret reference](#certificate-secret-reference).

###### Rocket example

To configure a Provider for Rocket, create a Secret with [the `address`](#address-example)
set to the Rocket [webhook URL](https://docs.rocket.chat/guides/administration/admin-panel/integrations#incoming-webhook-script),
and a `rocket` Provider with a [Secret reference](#secret-reference).

```yaml
---
apiVersion: notification.toolkit.fluxcd.io/v1beta3
kind: Provider
metadata:
  name: rocket
  namespace: default
spec:
  type: rocket
  secretRef:
    name: rocket-webhook
```

##### Google Chat

When `.spec.type` is set to `googlechat`, the controller will send a payload for
an [Event](events.md#event-structure) to the provided Google Chat [Address](#address).

The Event will be formatted into a [Google Chat card message](https://developers.google.com/chat/api/reference/rest/v1/cards-v1),
with the metadata added as a list of [key-value pairs](https://developers.google.com/chat/api/reference/rest/v1/cards-v1#keyvalue)
in a [widget](https://developers.google.com/chat/api/reference/rest/v1/cards-v1#widgetmarkup).

This Provider type does support the configuration of a [proxy URL](#https-proxy).

###### Google Chat example

To configure a Provider for Google Chat, create a Secret with [the `address`](#address-example)
set to the Google Chat [webhook URL](https://developers.google.com/chat/how-tos/webhooks#create_a_webhook),
and a `googlechat` Provider with a [Secret reference](#secret-reference).

```yaml
---
apiVersion: notification.toolkit.fluxcd.io/v1beta3
kind: Provider
metadata:
  name: google
  namespace: default
spec:
  type: googlechat
  secretRef:
    name: google-webhook
---
apiVersion: v1
kind: Secret
metadata:
  name: google-webhook
  namespace: default
stringData:
  address: https://chat.googleapis.com/v1/spaces/...
```

##### Google Pub/Sub

When `.spec.type` is set to `googlepubsub`, the controller will publish the payload of
an [Event](events.md#event-structure) on the Google Pub/Sub Topic ID provided in the
[Channel](#channel) field, which must exist in the GCP Project ID provided in the
[Address](#address) field.

This Provider type can optionally use the [Secret reference](#secret-reference) to
authenticate on the Google Pub/Sub API by using [JSON credentials](https://cloud.google.com/iam/docs/service-account-creds#key-types).
The credentials must be specified on [the `token`](#token-example) field of the Secret.

If no JSON credentials are specified, then the automatic authentication methods of
the Google libraries will take place, and therefore methods like Workload Identity
will be automatically attempted.

The Google identity effectively used for publishing messages must have
[the required permissions](https://cloud.google.com/iam/docs/understanding-roles#pubsub.publisher)
on the Pub/Sub Topic.

You can optionally add [attributes](https://cloud.google.com/pubsub/docs/samples/pubsub-publish-custom-attributes#pubsub_publish_custom_attributes-go)
to the Pub/Sub message using a [`headers` key in the referenced Secret](#http-headers-example).

This Provider type does not support the configuration of a [proxy URL](#https-proxy)
or [certificate secret reference](#certificate-secret-reference).

###### Google Pub/Sub with JSON Credentials and Custom Headers Example

To configure a Provider for Google Pub/Sub authenticating with JSON credentials and
custom headers, create a Secret with [the `token`](#token-example) set to the
necessary JSON credentials, [the `headers`](#http-headers-example) field set to a
YAML string-to-string dictionary, and a `googlepubsub` Provider with the associated
[Secret reference](#secret-reference).

```yaml
---
apiVersion: notification.toolkit.fluxcd.io/v1beta3
kind: Provider
metadata:
  name: googlepubsub-provider
  namespace: desired-namespace
spec:
  type: googlepubsub
  address: <GCP Project ID>
  channel: <Pub/Sub Topic ID>
  secretRef:
    name: googlepubsub-provider-creds
---
apiVersion: v1
kind: Secret
metadata:
  name: googlepubsub-provider-creds
  namespace: desired-namespace
stringData:
  token: <GCP JSON credentials>
  headers: |
    attr1-name: attr1-value
    attr2-name: attr2-value
```

##### Opsgenie

When `.spec.type` is set to `opsgenie`, the controller will send a payload for
an [Event](events.md#event-structure) to the provided Opsgenie [Address](#address).

The Event will be formatted into a [Opsgenie alert](https://docs.opsgenie.com/docs/alert-api#section-create-alert-request),
with the metadata added to the [`details` field](https://docs.opsgenie.com/docs/alert-api#create-alert)
as a list of key-value pairs.

This Provider type does support the configuration of a [proxy URL](#https-proxy)
and [certificate secret reference](#certificate-secret-reference).

###### Opsgenie example

To configure a Provider for Opsgenie, create a Secret with [the `token`](#token-example)
set to the Opsgenie [API token](https://support.atlassian.com/opsgenie/docs/create-a-default-api-integration/),
and a `opsgenie` Provider with a [Secret reference](#secret-reference) and the
`address` set to `https://api.opsgenie.com/v2/alerts`.

```yaml
---
apiVersion: notification.toolkit.fluxcd.io/v1beta3
kind: Provider
metadata:
  name: opsgenie
  namespace: default
spec:
  type: opsgenie
  address: https://api.opsgenie.com/v2/alerts
  secretRef:
    name: opsgenie-token
---
apiVersion: v1
kind: Secret
metadata:
  name: opsgenie-token
  namespace: default
stringData:
    token: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

##### PagerDuty

When `.spec.type` is set to `pagerduty`, the controller will send a payload for
an [Event](events.md#event-structure) to the provided PagerDuty [Address](#address).

The Event will be formatted into an [Event API v2](https://developer.pagerduty.com/api-reference/368ae3d938c9e-send-an-event-to-pager-duty) payload,
triggering or resolving an incident depending on the event's `Severity`.

The provider will also send [Change Events](https://developer.pagerduty.com/api-reference/95db350959c37-send-change-events-to-the-pager-duty-events-api)
for `info` level `Severity`, which will be displayed in the PagerDuty service's timeline to track changes.

This Provider type supports the configuration of a [proxy URL](#https-proxy)
and [certificate secret reference](#certificate-secret-reference).

The [Channel](#channel) is used to set the routing key to send the event to the appropriate integration.

###### PagerDuty example

To configure a Provider for Pagerduty, create a `pagerduty` Provider,
set `address` to the integration URL and `channel` set to
the integration key (also known as a routing key) for your [service](https://support.pagerduty.com/docs/services-and-integrations#create-a-generic-events-api-integration)
or [event orchestration](https://support.pagerduty.com/docs/event-orchestration).

When adding an integration for a service on PagerDuty, it is recommended to use `Events API v2` integration.

**Note**: PagerDuty does not support Change Events when sent to global integrations, such as event orchestration.

```yaml
---
apiVersion: notification.toolkit.fluxcd.io/v1beta3
kind: Provider
metadata:
  name: pagerduty
  namespace: default
spec:
  type: pagerduty
  address: https://events.pagerduty.com
  channel: <integrationKey>
```
If you are sending to a service integration, it is recommended to set your Alert to filter to
only those sources you want to trigger an incident for that service. For example:

```yaml
---
apiVersion: notification.toolkit.fluxcd.io/v1beta3
kind: Alert
metadata:
  name: my-service-pagerduty
  namespace: default
spec:
  providerRef:
    name: pagerduty
  eventSources:
    - kind: HelmRelease
      name: my-service
      namespace: default
```

##### Prometheus Alertmanager

When `.spec.type` is set to `alertmanager`, the controller will send a payload for
an [Event](events.md#event-structure) to the provided Prometheus Alertmanager
[Address](#address).

The Event will be formatted into a `firing` [Prometheus Alertmanager
alert](https://prometheus.io/docs/alerting/latest/notifications/#alert), with
the metadata added to the `labels` fields, and the `message` (and optional
`.metadata.summary`) added as annotations. Event timestamp will be used to set
alert start time (`.StartsAt`).

In addition to the metadata from the Event, the following labels will be added:

| Label     | Description                                                                                          |
|-----------|------------------------------------------------------------------------------------------------------|
| alertname | The string Flux followed by the Kind and the reason for the event e.g `FluxKustomizationProgressing` |
| severity  | The severity of the event (`error` or `info`)                                                        |
| reason    | The machine readable reason for the objects transition into the current status                       |
| kind      | The kind of the involved object associated with the event                                            |
| name      | The name of the involved object associated with the event                                            |
| namespace | The namespace of the involved object associated with the event                                       |

Note that due to the way other Flux controllers currently emit events, there's
no way for notification-controller to figure out the time the event ends to set
`.EndsAt` (a reasonable estimate being double the reconciliation interval of the
resource involved) that doesn't involve a Kubernetes API roundtrip. A
possible workaround could be setting
[`global.resolve_timeout`][am_config_global] to an interval large enough for
events to reoccur:

[am_config_global]: https://prometheus.io/docs/alerting/latest/configuration/#file-layout-and-global-settings

```yaml
global:
  resolve_timeout: 1h
```

This Provider type does support the configuration of a [proxy URL](#https-proxy)
and [certificate secret reference](#certificate-secret-reference).

###### Prometheus Alertmanager example

To configure a Provider for Prometheus Alertmanager, authentication can be done using either Basic Authentication or a Bearer Token. 
Both methods are supported, but using authentication is optional based on your setup.

Basic Authentication: 
Create a Secret with [the `username` and `password`](#secret-reference) set to the Basic Auth
credentials, and the [`.spec.address`](#address) field set to the Prometheus Alertmanager
[HTTP API URL](https://prometheus.io/docs/alerting/latest/https/#http-traffic).

```yaml
---
apiVersion: notification.toolkit.fluxcd.io/v1beta3
kind: Provider
metadata:
  name: alertmanager
  namespace: default
spec:
  type: alertmanager
  address: https://<alertmanager-hostport>/api/v2/alerts/
  secretRef:
    name: alertmanager-basic-auth
---
apiVersion: v1
kind: Secret
metadata:
  name: alertmanager-basic-auth
  namespace: default
stringData:
  username: <username>
  password: <password>
```

Bearer Token Authentication:
Create a Secret with [the `token`](#token-example), and an `alertmanager` Provider with a [Secret
reference](#secret-reference) and the Prometheus Alertmanager [HTTP API
URL](https://prometheus.io/docs/alerting/latest/https/#http-traffic) set directly in the `.spec.address` field.

```yaml
---
apiVersion: notification.toolkit.fluxcd.io/v1beta3
kind: Provider
metadata:
  name: alertmanager
  namespace: default
spec:
  type: alertmanager
  address: https://<alertmanager-hostport>/api/v2/alerts/
  secretRef:
    name: alertmanager-token
---
apiVersion: v1
kind: Secret
metadata:
  name: alertmanager-token
  namespace: default
stringData:
  token: <token>
```

##### Webex

When `.spec.type` is set to `webex`, the controller will send a payload for
an [Event](events.md#event-structure) to the provided Webex [Address](#address).

The Event will be formatted into a message string, with the metadata attached
as a list of key-value pairs, and send as a [Webex message](https://developer.webex.com/docs/api/v1/messages/create-a-message).

The [Channel](#channel) is used to set the ID of the room to send the message
to.

This Provider type does support the configuration of a [proxy URL](#https-proxy)
and [certificate secret reference](#certificate-secret-reference).

###### Webex example

To configure a Provider for Webex, create a Secret with [the `token`](#token-example)
set to the Webex [access token](https://developer.webex.com/docs/api/getting-started#authentication),
and a `webex` Provider with a [Secret reference](#secret-reference) and the
`address` set to `https://webexapis.com/v1/messages`.

**Note:** To be able to send messages to a Webex room, the bot needs to be
added to the room. Failing to do so will result in 404 errors, logged by the
controller.

```yaml
---
apiVersion: notification.toolkit.fluxcd.io/v1beta3
kind: Provider
metadata:
  name: webex
  namespace: default
spec:
  type: webex
  address: https://webexapis.com/v1/messages
  channel: <webexSpaceRoomID>
  secretRef:
    name: webex-token
---
apiVersion: v1
kind: Secret
metadata:
  name: webex-token
  namespace: default
stringData:
  token: <bot-token>
```

##### NATS

When `.spec.type` is set to `nats`, the controller will publish the payload of
an [Event](events.md#event-structure) on the [NATS Subject](https://docs.nats.io/nats-concepts/subjects) provided in the
[Channel](#channel) field, using the server specified in the [Address](#address) field.

This Provider type can optionally use the [Secret reference](#secret-reference) to
authenticate to the NATS server using [Username/Password](https://docs.nats.io/using-nats/developer/connecting/userpass).
The credentials must be specified in [the `username`](#username-example) and `password` fields of the Secret.
Alternatively, NATS also supports passing the credentials with [the server URL](https://docs.nats.io/using-nats/developer/connecting/userpass#connecting-with-a-user-password-in-the-url). In this case the `address` should be provided through a 
Secret reference.

Additionally if using credentials, the User must have [authorization](https://docs.nats.io/running-a-nats-service/configuration/securing_nats/authorization) to publish on the Subject provided.

###### NATS with Username/Password Credentials Example

To configure a Provider for NATS authenticating with Username/Password, create a Secret with the
`username` and `password` fields set, and add a `nats` Provider with the associated
[Secret reference](#secret-reference).

```yaml
---
apiVersion: notification.toolkit.fluxcd.io/v1beta3
kind: Provider
metadata:
  name: nats-provider
  namespace: desired-namespace
spec:
  type: nats
  address: <NATS Server URL>
  channel: <Subject>
  secretRef:
    name: nats-provider-creds
---
apiVersion: v1
kind: Secret
metadata:
  name: nats-provider-creds
  namespace: desired-namespace
stringData:
  username: <NATS Username>
  password: <NATS Password>
```

### Address

`.spec.address` is an optional field that specifies the endpoint where the events are posted.
The meaning of endpoint here depends on the specific Provider type being used.
For the `generic` Provider for example this is an HTTP/S address.
For other Provider types this could be a project ID or a namespace.

If the address contains sensitive information such as tokens or passwords, it is 
recommended to store the address in the Kubernetes secret referenced by `.spec.secretRef.name`.
When the referenced Secret contains an `address` key, the `.spec.address` value is ignored.

### Channel

`.spec.channel` is an optional field that specifies the channel where the events are posted.

### Username

`.spec.username` is an optional field that specifies the username used to post
the events. Can be overwritten with a [Secret reference](#secret-reference).

### Secret reference

`.spec.secretRef.name` is an optional field to specify a name reference to a
Secret in the same namespace as the Provider, containing the authentication
credentials for the provider API.

The Kubernetes secret can have any of the following keys:

- `address` - overrides `.spec.address`
- `proxy` - overrides `.spec.proxy` (deprecated, use `.spec.proxySecretRef` instead. **Support for this key will be removed in v1**)
- `token` - used for authentication
- `username` - overrides `.spec.username`
- `password` - used for authentication, often in combination with `username` (or `.spec.username`)
- `headers` - HTTP headers values included in the POST request

#### Address example

For providers which embed tokens or other sensitive information in the URL,
the incoming webhooks address can be stored in the secret using the `address` key:

```yaml
---
apiVersion: v1
kind: Secret
metadata:
  name: my-provider-url
  namespace: default
stringData:
  address: "https://webhook.example.com/token"
```

#### Token example

For providers which require token based authentication, the API token
can be stored in the secret using the `token` key:

```yaml
---
apiVersion: v1
kind: Secret
metadata:
  name: my-provider-auth
  namespace: default
stringData:
  token: "my-api-token"
```

#### HTTP headers example

For providers which require specific HTTP headers to be attached to the POST request,
the headers can be set in the secret using the `headers` key:

```yaml
---
apiVersion: v1
kind: Secret
metadata:
  name: my-provider-headers
  namespace: default
stringData:
  headers: |
    Authorization: my-api-token
    X-Forwarded-Proto: https
```

#### Proxy auth example

Some networks need to use an authenticated proxy to access external services.
The recommended approach is to use `.spec.proxySecretRef` with a dedicated Secret:

```yaml
---
apiVersion: v1
kind: Secret
metadata:
  name: my-provider-proxy
  namespace: default
stringData:
  address: "http://proxy_url:proxy_port"
  username: "proxy_username"
  password: "proxy_password"
```

**Legacy approach (deprecated):**
The proxy address can also be stored in the main secret to hide parameters like the username and password:

```yaml
---
apiVersion: v1
kind: Secret
metadata:
  name: my-provider-proxy-legacy
  namespace: default
stringData:
  proxy: "http://username:password@proxy_url:proxy_port"
```

### Certificate secret reference

`.spec.certSecretRef` is an optional field to specify a name reference to a
Secret in the same namespace as the Provider, containing TLS certificates for
secure communication. The secret must be of type `kubernetes.io/tls` or `Opaque`.

#### Supported configurations

- **Mutual TLS (mTLS)**: Client certificate authentication (provide `tls.crt` + `tls.key`, optionally with `ca.crt`)
- **CA-only**: Server authentication (provide `ca.crt` only)

#### Mutual TLS Authentication

Mutual TLS authentication allows for secure client-server communication using
client certificates stored in Kubernetes secrets. Both `tls.crt` and `tls.key`
must be specified together for client certificate authentication. The `ca.crt`
field is optional but required when connecting to servers with self-signed certificates.

##### Providers supporting client certificate authentication

The following providers support client certificate authentication:

| Provider Type        | Description                    |
|---------------------|--------------------------------|
| `alertmanager`      | Prometheus Alertmanager        |
| `azuredevops`       | Azure DevOps                   |
| `bitbucket`         | Bitbucket                      |
| `bitbucketserver`   | BitBucket Server/Data Center   |
| `datadog`           | DataDog                        |
| `discord`           | Discord webhooks               |
| `forwarder`         | Generic forwarder              |
| `gitea`             | Gitea                          |
| `github`            | GitHub                         |
| `githubdispatch`    | GitHub Dispatch                |
| `gitlab`            | GitLab                         |
| `grafana`           | Grafana annotations API        |
| `matrix`            | Matrix rooms                   |
| `msteams`           | Microsoft Teams                |
| `opsgenie`          | Opsgenie alerts                |
| `pagerduty`         | PagerDuty events               |
| `rocket`            | Rocket.Chat                    |
| `sentry`            | Sentry                         |
| `slack`             | Slack API                      |
| `webex`             | Webex messages                 |

Support for client certificate authentication is being expanded to additional providers over time.

##### Example: mTLS Configuration

```yaml
---
apiVersion: notification.toolkit.fluxcd.io/v1beta3
kind: Provider
metadata:
  name: my-webhook-mtls
  namespace: default
spec:
  type: generic
  address: https://my-webhook.internal
  certSecretRef:
    name: my-mtls-certs
---
apiVersion: v1
kind: Secret
metadata:
  name: my-mtls-certs
  namespace: default
type: kubernetes.io/tls # or Opaque
stringData:
  tls.crt: |
    -----BEGIN CERTIFICATE-----
    <client certificate>
    -----END CERTIFICATE-----
  tls.key: |
    -----BEGIN PRIVATE KEY-----
    <client private key>
    -----END PRIVATE KEY-----
  ca.crt: |
    -----BEGIN CERTIFICATE-----
    <certificate authority certificate>
    -----END CERTIFICATE-----
```

#### CA Certificate Authentication

CA certificate authentication provides server authentication when connecting to
HTTPS endpoints with self-signed or custom CA certificates. Only the `ca.crt`
field is required for this configuration.

##### Example: CA Certificate Configuration

```yaml
---
apiVersion: notification.toolkit.fluxcd.io/v1beta3
kind: Provider
metadata:
  name: my-webhook
  namespace: flagger-system
spec:
  type: generic
  address: https://my-webhook.internal
  certSecretRef:
    name: my-ca-crt
---
apiVersion: v1
kind: Secret
metadata:
  name: my-ca-crt
  namespace: default
type: kubernetes.io/tls # or Opaque
stringData:
  ca.crt: |
    -----BEGIN CERTIFICATE-----
    <certificate authority certificate>
    -----END CERTIFICATE-----
```

**Warning:** Support for the `caFile` key has been
deprecated. If you have any Secrets using this key,
the controller will log a deprecation warning.

### HTTP/S proxy

`.spec.proxy` is an optional field to specify an HTTP/S proxy address.
**Warning:** This field is deprecated, use `.spec.proxySecretRef` instead. **Support for this field will be removed in v1.**

`.spec.proxySecretRef` is an optional field to specify a name reference to a
Secret in the same namespace as the Provider, containing the proxy configuration.
The Secret should contain an `address` key with the HTTP/S address of the proxy server.
Optional `username` and `password` keys can be provided for proxy authentication.

If the proxy address contains sensitive information such as basic auth credentials, it is
recommended to use `.spec.proxySecretRef` instead of `.spec.proxy`.
When `.spec.proxySecretRef` is specified, both `.spec.proxy` and the `proxy` key from
`.spec.secretRef` are ignored.

### Timeout

`.spec.timeout` is an optional field to specify the timeout for the
HTTP/S request sent to the provider endpoint.
The value must be in a
[Go recognized duration string format](https://pkg.go.dev/time#ParseDuration),
e.g. `5m30s` for a timeout of five minutes and thirty seconds.

### Suspend

`.spec.suspend` is an optional field to suspend the provider.
When set to `true`, the controller will stop sending events to this provider.
When the field is set to `false` or removed, it will resume.

## Working with Providers


### Grafana

To send notifications to [Grafana annotations API](https://grafana.com/docs/grafana/latest/http_api/annotations/),
enable the annotations on a Dashboard like so:

- Annotations > Query > Enable Match any
- Annotations > Query > Tags (Add Tag: `flux`)

If Grafana has authentication configured, create a Kubernetes Secret with the API token:

```shell
kubectl create secret generic grafana-token \
--from-literal=token=<grafana-api-key> \
```

Grafana can also use basic authorization to authenticate the requests, if both the token and
the username/password are set in the secret, then token takes precedence over`basic auth:

```shell
kubectl create secret generic grafana-token \
--from-literal=username=<your-grafana-username> \
--from-literal=password=<your-grafana-password>
```

Create a provider of type `grafana` and reference the `grafana-token` secret:

```yaml
apiVersion: notification.toolkit.fluxcd.io/v1beta3
kind: Provider
metadata:
  name: grafana
  namespace: default
spec:
  type: grafana
  address: https://<grafana-url>/api/annotations
  secretRef:
    name: grafana-token
```

Besides the tag `flux` and the tag containing the reporting controller (e.g. `source-controller`),
the event metadata is also included as tags of the form `${metadataKey}: ${metadataValue}`, and
the tags `kind: ${event.InvolvedObject.Kind}`, `name: ${event.InvolvedObject.Name}` and
`namespace: ${event.InvolvedObject.Namespace}` are also included.

### GitHub dispatch

The `githubdispatch` provider generates GitHub events of type
[`repository_dispatch`](https://docs.github.com/en/rest/reference/repos#create-a-repository-dispatch-event)
for the selected repository. The `repository_dispatch` events can be used to trigger GitHub Actions workflow.

The request includes the `event_type` and `client_payload` fields:

- `event_type` is generated from the involved object in the format `{Kind}/{Name}.{Namespace}`.
- `client_payload` contains the [Flux event](events.md).

### Setting up the GitHub dispatch provider

```yaml
---
apiVersion: notification.toolkit.fluxcd.io/v1beta3
kind: Provider
metadata:
  name: github-dispatch
  namespace: flux-system
spec:
  type: githubdispatch
  address: https://github.com/stefanprodan/podinfo
  secretRef:
    name: auth-secret
```

The `address` is the address of your repository where you want to send webhooks to trigger GitHub workflows.

GitHub uses [personal access tokens](https://docs.github.com/en/free-pro-team@latest/github/authenticating-to-github/creating-a-personal-access-token)
or [GitHub app](https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/authenticating-as-a-github-app-installation)
for authentication with its API:

#### GitHub personal access token

The provider requires a secret in the same format, with the personal access token as the value for the token key:

```yaml
---
apiVersion: v1
kind: Secret
metadata:
  name: auth-secret
  namespace: default
stringData:
  token: <personal-access-tokens>
```

#### GitHub App

To use Github App authentication, make sure the GitHub App is registered and
installed with the necessary permissions and the github app secret is created as
described
[here](https://fluxcd.io/flux/components/source/gitrepositories/#github).

#### Setting up a GitHub workflow

To trigger a GitHub Actions workflow when a Flux Kustomization finishes reconciling,
you need to set the event type for the repository_dispatch trigger to match the Flux object ID:

```yaml
name: test-github-dispatch-provider
on:
  repository_dispatch:
    types: [Kustomization/podinfo.flux-system]
```

Assuming that we deploy all Flux kustomization resources in the same namespace,
it will be useful to have a unique kustomization resource name for each application.
This will allow you to use only `event_type` to trigger tests for the exact application.

Let's say we have following folder structure for applications kustomization manifests:

```bash
apps/
├── app1
│   └── overlays
│       ├── production
│       └── staging
└── app2
    └── overlays
        ├── production
        └── staging
```

You can then create a flux kustomization resource for the app to have unique `event_type` per app.
The kustomization manifest for app1/staging:

```yaml
apiVersion: kustomize.toolkit.fluxcd.io/v1beta3
kind: Kustomization
metadata:
  name: app1
  namespace: flux-system
spec:
  path: "./app1/staging"
```

You would also like to know from the notification which cluster is being used for deployment.
You can add the `spec.summary` field to the Flux alert configuration to mention the relevant cluster:

```yaml
apiVersion: notification.toolkit.fluxcd.io/v1beta3
kind: Alert
metadata:
  name: github-dispatch
  namespace: flux-system
spec:
  summary: "staging (us-west-2)"
  providerRef:
    name: github-dispatch
  eventSeverity: info
  eventSources:
    - kind: Kustomization
      name: 'podinfo'
```

Now you can trigger the tests in the GitHub workflow for app1 in a staging cluster when
the app1 resources defined in `./app1/staging/` are reconciled by Flux:

```yaml
name: test-github-dispatch-provider
on:
  repository_dispatch:
    types: [Kustomization/podinfo.flux-system]
jobs:
  run-tests-staging:
    if: github.event.client_payload.metadata.summary == 'staging (us-west-2)'
    runs-on: ubuntu-18.04
    steps:
    - name: Run tests
      run: echo "running tests.."
```

### Azure Event Hub

The Azure Event Hub provider supports the following authentication methods, 
- [Managed
  Identity](https://learn.microsoft.com/en-us/azure/event-hubs/authenticate-managed-identity)
- [SAS](https://docs.microsoft.com/en-us/azure/event-hubs/authorize-access-shared-access-signature)
  based.
- [JWT](https://docs.microsoft.com/en-us/azure/event-hubs/authenticate-application) (Deprecated)

#### Managed Identity

Managed identity authentication can be setup using Azure Workload identity. 

##### Pre-requisites

- Ensure Workload Identity is properly [set up on your
  cluster](https://learn.microsoft.com/en-us/azure/aks/workload-identity-deploy-cluster#create-an-aks-cluster).

##### Configure workload identity

- Create a managed identity to access Azure Event Hub. 
- Grant the managed identity the necessary permissions to send events to Azure
  Event hub as described
  [here](https://learn.microsoft.com/en-us/azure/event-hubs/authenticate-managed-identity#to-assign-azure-roles-using-the-azure-portal).

- Establish a federated identity credential between the managed identity and the
  service account to be used for authentication. Ensure the federated credential
  uses the correct namespace and name of the service account. For more details,
  please refer to this
  [guide](https://azure.github.io/azure-workload-identity/docs/quick-start.html#6-establish-federated-identity-credential-between-the-identity-and-the-service-account-issuer--subject).

##### Single tenant approach

This approach uses the notification-controller service account for setting up
authentication.

- In the default installation, the notification-controller service account is
  located in the `flux-system` namespace with name `notification-controller`.

- Configure workload identity with notification-controller as described in the
  docs [here](/flux/installation/configuration/workload-identity/).

##### Multi-tenant approach

For multi-tenant clusters, set `.spec.serviceAccountName` of the provider to
the service account to be used for authentication. Ensure that the service
account has the
[annotations](https://learn.microsoft.com/en-us/azure/aks/workload-identity-overview?tabs=dotnet#service-account-annotations)
for the client-id and tenant-id of the managed identity.

For a complete guide on how to set up authentication for an Azure Event Hub,
see the integration [docs](/flux/integrations/azure/).

#### SAS based auth

When using SAS auth, we only use the `address` field in the secret.

```yaml
---
apiVersion: notification.toolkit.fluxcd.io/v1beta3
kind: Provider
metadata:
  name: azure
  namespace: default
spec:
  type: azureeventhub
  secretRef:
    name: azure-webhook
---
apiVersion: v1
kind: Secret
metadata:
  name: azure-webhook
  namespace: default
stringData:
  address: <SAS-URL>
```

Assuming that you have created the Azure event hub and namespace you should be
able to use a similar command to get your connection string. This will give you
the default Root SAS, which is NOT supposed to be used in production.

```shell
az eventhubs namespace authorization-rule keys list --resource-group <rg-name> --namespace-name <namespace-name> --name RootManageSharedAccessKey -o tsv --query primaryConnectionString
# The output should look something like this:
Endpoint=sb://fluxv2.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=yoursaskeygeneatedbyazure;EntityPath=youreventhub
```

To create the needed secret:

```shell
kubectl create secret generic azure-webhook \
--from-literal=address="Endpoint=sb://fluxv2.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=yoursaskeygeneatedbyazure"
```

#### JWT based auth (Deprecated)

In JWT we use 3 input values. Channel, token and address. We perform the
following translation to match we the data we need to communicate with Azure
Event Hub.

- channel = Azure Event Hub namespace
- address = Azure Event Hub name 
- token   = JWT

```yaml
---
apiVersion: notification.toolkit.fluxcd.io/v1beta3
kind: Provider
metadata:
  name: azure
  namespace: default
spec:
  type: azureeventhub
  address: <event-hub-name>
  channel: <event-hub-namespace>
  secretRef:
    name: azure-token
---
apiVersion: v1
kind: Secret
metadata:
  name: azure-token
  namespace: default
stringData:
  token: <event-hub-token>
```

The controller doesn't take any responsibility for the JWT token to be updated.
You need to use a secondary tool to make sure that the token in the secret is
renewed.

If you want to make a easy test assuming that you have setup a Azure Enterprise
application and you called it event-hub you can follow most of the bellow
commands. You will need to provide the `client_secret` that you got when
generating the Azure Enterprise Application.

```shell
export AZURE_CLIENT=$(az ad app list --filter "startswith(displayName,'event-hub')" --query '[].appId' |jq -r '.[0]')
export AZURE_SECRET='secret-client-secret-generated-at-creation'
export AZURE_TENANT=$(az account show -o tsv --query tenantId)

curl -X GET --data 'grant_type=client_credentials' --data "client_id=$AZURE_CLIENT" --data "client_secret=$AZURE_SECRET" --data 'resource=https://eventhubs.azure.net' -H 'Content-Type: application/x-www-form-urlencoded' https://login.microsoftonline.com/$AZURE_TENANT/oauth2/token |jq .access_token
```

Use the output you got from `curl` and add it to your secret like bellow:

```shell
kubectl create secret generic azure-token \
--from-literal=token='A-valid-JWT-token'
```

### Git Commit Status Updates

The notification-controller can mark Git commits as reconciled by posting
Flux `Kustomization` events to the origin repository using Git SaaS providers APIs.

#### Example

The following is an example of how to update the Git commit status for the GitHub repository where
Flux was bootstrapped with `flux bootstrap github --owner=my-gh-org --repository=my-gh-repo`.

```yaml
apiVersion: notification.toolkit.fluxcd.io/v1beta3
kind: Provider
metadata:
  name: github-status
  namespace: flux-system
spec:
  type: github
  address: https://github.com/my-gh-org/my-gh-repo
  secretRef:
    name: github-token
---
apiVersion: notification.toolkit.fluxcd.io/v1beta3
kind: Alert
metadata:
  name: github-status
  namespace: flux-system
spec:
  providerRef:
    name: github-status
  eventSources:
    - kind: Kustomization
      name: flux-system
```

#### GitHub

When `.spec.type` is set to `github`, the referenced secret can contain a
personal access token OR github app details. 

##### Personal Access Token

To use personal access tokens, the secret must contain a key called `token` with
the value set to a [GitHub personal access
token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token).
The token must have permissions to update the commit status for the GitHub
repository specified in `.spec.address`.

You can create the secret with `kubectl` like this:

```shell
kubectl create secret generic github-token --from-literal=token=<GITHUB-TOKEN>
```

##### GitHub App

To use [Github App
authentication](https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/authenticating-as-a-github-app-installation),
make sure the GitHub App is registered and installed with the necessary
permissions to update the commit status and the github app secret is created as
described
[here](https://fluxcd.io/flux/components/source/gitrepositories/#github).

#### GitLab

When `.spec.type` is set to `gitlab`, the referenced secret must contain a key called `token` with the value set to a
[GitLab personal access token](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html). If available group and project access tokens can also be used.

The token must have permissions to update the commit status for the GitLab repository specified in `.spec.address` (grant it the `api` scope).

You can create the secret with `kubectl` like this:

```shell
kubectl create secret generic gitlab-token --from-literal=token=<GITLAB-TOKEN>
```

For gitlab.com and current self-hosted gitlab installations that support only the Gitlab v4 API you need to use the project ID in the address instead of the project name.
Use an address like `https://gitlab.com/1234` (with `1234` being the project ID).
You can find out the ID by opening the project in the browser and clicking on the three dot button in the top right corner. The menu that opens shows the project ID.

#### Gitea

When `.spec.type` is set to `gitea`, the referenced secret must contain a key called `token` with the value set to a
[Gitea token](https://docs.gitea.io/en-us/api-usage/#generating-and-listing-api-tokens).

The token must have at least the `write:repository` permission for the provider to
update the commit status for the Gitea repository specified in `.spec.address`.

{{% alert color="info" title="Gitea 1.20.0 & 1.20.1" %}}
Due to a bug in Gitea 1.20.0 and 1.20.1, these versions require the additional
`read:misc` scope to be applied to the token.
{{% /alert %}}

You can create the secret with `kubectl` like this:

```shell
kubectl create secret generic gitea-token --from-literal=token=<GITEA-TOKEN>
```

#### BitBucket

When `.spec.type` is set to `bitbucket`, the referenced secret must contain a key called `token` with the value
set to a BitBucket username and an
[app password](https://support.atlassian.com/bitbucket-cloud/docs/app-passwords/#Create-an-app-password)
in the format `<username>:<app-password>`.

The app password must have `Repositories (Read/Write)` permission for
the BitBucket repository specified in `.spec.address`.

You can create the secret with `kubectl` like this:

```shell
kubectl create secret generic bitbucket-token --from-literal=token=<username>:<app-password>
```

#### BitBucket Server/Data Center

When `.spec.type` is set to `bitbucketserver`, the following auth methods are available:

- Basic Authentication (username/password)
- [HTTP access tokens](https://confluence.atlassian.com/bitbucketserver/http-access-tokens-939515499.html)

For Basic Authentication, the referenced secret must contain a `password` field. The `username` field can either come from the [`.spec.username` field of the Provider](https://fluxcd.io/flux/components/notification/providers/#username) or can be defined in the referenced secret.

You can create the secret with `kubectl` like this:

```shell
kubectl create secret generic bb-server-username-password --from-literal=username=<username> --from-literal=password=<password>
```

For HTTP access tokens, the secret can be created with `kubectl` like this:

```shell
kubectl create secret generic bb-server-token --from-literal=token=<token>
```

The HTTP access token must have `Repositories (Read/Write)` permission for
the repository specified in `.spec.address`.

**NOTE:** Please provide HTTPS clone URL in the `address` field of this provider. SSH URLs are not supported by this provider type.

#### Azure DevOps

The following authentication methods can be used when `.spec.type` is set to
`azuredevops`. 

- [Managed Identity](https://learn.microsoft.com/en-us/azure/event-hubs/authenticate-managed-identity)
- [Personal Access Token](#pat)

#### Managed Identity

Managed Identity authentication can be setup using Azure Workload identity. 

##### Pre-requisites

- Ensure Workload Identity is properly 
  [set up on your cluster](https://learn.microsoft.com/en-us/azure/aks/workload-identity-deploy-cluster#create-an-aks-cluster).

##### Configure Workload Identity

- Create a Managed Identity and grant the necessary permissions to list/update
  commit status.
- Establish a federated identity credential between the managed identity and the
  service account to be used for authentication. Ensure the federated credential
  uses the correct namespace and name of the service account. For more details,
  please refer to this
  [guide](https://azure.github.io/azure-workload-identity/docs/quick-start.html#6-establish-federated-identity-credential-between-the-identity-and-the-service-account-issuer--subject).
The service account used for authentication can be single-tenant
(controller-level) or multi-tenant(object-level). For a complete guide on how to
set up authentication, see the integration [docs](/flux/integrations/azure/).

#### PAT

The `.spec.secretRef` must contain a key called `token` with the value set to a
[Azure DevOps personal access token](https://docs.microsoft.com/en-us/azure/devops/organizations/accounts/use-personal-access-tokens-to-authenticate?view=azure-devops&tabs=preview-page).

The token must have permissions to update the commit status for the Azure DevOps repository specified in `.spec.address`.

You can create the secret with `kubectl` like this:

```shell
kubectl create secret generic azuredevops-token --from-literal=token=<AZURE-TOKEN>
```

#### Custom Commit Status Messages

Git providers supporting commit status updates can use a CEL expression to build a custom commit status message by setting the optional field `spec.commitStatusExpr`:

```yaml
apiVersion: notification.toolkit.fluxcd.io/v1beta3
kind: Provider
metadata:
  name: github-status
  namespace: flux-system
spec:
  type: github
  address: https://github.com/my-gh-org/my-gh-repo
  secretRef:
    name: github-token
  commitStatusExpr: "(event.involvedObject.kind + '/' + event.involvedObject.name + '/' + event.metadata.foo + '/' + provider.metadata.uid.split('-').first().value()).lowerAscii()"
```

The CEL expression can access the following variables:
- `event`: The Flux event object containing metadata about the reconciliation
- `provider`: The Provider object
- `alert`: The Alert object

If the `spec.commitStatusExpr` field is not specified, the notification-controller will use a default commit status message based on the involved object kind, name, and a truncated provider UID to generate a commit status (e.g. `kustomization/gitops-system/0c9c2e41`).

A useful tool for building and testing CEL expressions is the [CEL Playground](https://playcel.undistro.io/).
