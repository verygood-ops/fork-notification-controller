<h1>Notification API reference v1beta3</h1>
<p>Packages:</p>
<ul class="simple">
<li>
<a href="#notification.toolkit.fluxcd.io%2fv1beta3">notification.toolkit.fluxcd.io/v1beta3</a>
</li>
</ul>
<h2 id="notification.toolkit.fluxcd.io/v1beta3">notification.toolkit.fluxcd.io/v1beta3</h2>
<p>Package v1beta3 contains API Schema definitions for the notification v1beta3 API group.</p>
Resource Types:
<ul class="simple"><li>
<a href="#notification.toolkit.fluxcd.io/v1beta3.Alert">Alert</a>
</li><li>
<a href="#notification.toolkit.fluxcd.io/v1beta3.Provider">Provider</a>
</li></ul>
<h3 id="notification.toolkit.fluxcd.io/v1beta3.Alert">Alert
</h3>
<p>Alert is the Schema for the alerts API</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code><br>
string</td>
<td>
<code>notification.toolkit.fluxcd.io/v1beta3</code>
</td>
</tr>
<tr>
<td>
<code>kind</code><br>
string
</td>
<td>
<code>Alert</code>
</td>
</tr>
<tr>
<td>
<code>metadata</code><br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br>
<em>
<a href="#notification.toolkit.fluxcd.io/v1beta3.AlertSpec">
AlertSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>providerRef</code><br>
<em>
<a href="https://pkg.go.dev/github.com/fluxcd/pkg/apis/meta#LocalObjectReference">
github.com/fluxcd/pkg/apis/meta.LocalObjectReference
</a>
</em>
</td>
<td>
<p>ProviderRef specifies which Provider this Alert should use.</p>
</td>
</tr>
<tr>
<td>
<code>eventSeverity</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>EventSeverity specifies how to filter events based on severity.
If set to &lsquo;info&rsquo; no events will be filtered.</p>
</td>
</tr>
<tr>
<td>
<code>eventSources</code><br>
<em>
<a href="https://pkg.go.dev/github.com/fluxcd/notification-controller/api/v1#CrossNamespaceObjectReference">
[]github.com/fluxcd/notification-controller/api/v1.CrossNamespaceObjectReference
</a>
</em>
</td>
<td>
<p>EventSources specifies how to filter events based
on the involved object kind, name and namespace.</p>
</td>
</tr>
<tr>
<td>
<code>inclusionList</code><br>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>InclusionList specifies a list of Golang regular expressions
to be used for including messages.</p>
</td>
</tr>
<tr>
<td>
<code>eventMetadata</code><br>
<em>
map[string]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>EventMetadata is an optional field for adding metadata to events dispatched by the
controller. This can be used for enhancing the context of the event. If a field
would override one already present on the original event as generated by the emitter,
then the override doesn&rsquo;t happen, i.e. the original value is preserved, and an info
log is printed.</p>
</td>
</tr>
<tr>
<td>
<code>exclusionList</code><br>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ExclusionList specifies a list of Golang regular expressions
to be used for excluding messages.</p>
</td>
</tr>
<tr>
<td>
<code>summary</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Summary holds a short description of the impact and affected cluster.
Deprecated: Use EventMetadata instead.</p>
</td>
</tr>
<tr>
<td>
<code>suspend</code><br>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Suspend tells the controller to suspend subsequent
events handling for this Alert.</p>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="notification.toolkit.fluxcd.io/v1beta3.Provider">Provider
</h3>
<p>Provider is the Schema for the providers API</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code><br>
string</td>
<td>
<code>notification.toolkit.fluxcd.io/v1beta3</code>
</td>
</tr>
<tr>
<td>
<code>kind</code><br>
string
</td>
<td>
<code>Provider</code>
</td>
</tr>
<tr>
<td>
<code>metadata</code><br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br>
<em>
<a href="#notification.toolkit.fluxcd.io/v1beta3.ProviderSpec">
ProviderSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>type</code><br>
<em>
string
</em>
</td>
<td>
<p>Type specifies which Provider implementation to use.</p>
</td>
</tr>
<tr>
<td>
<code>interval</code><br>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Interval at which to reconcile the Provider with its Secret references.
Deprecated and not used in v1beta3.</p>
</td>
</tr>
<tr>
<td>
<code>channel</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Channel specifies the destination channel where events should be posted.</p>
</td>
</tr>
<tr>
<td>
<code>username</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Username specifies the name under which events are posted.</p>
</td>
</tr>
<tr>
<td>
<code>address</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Address specifies the endpoint, in a generic sense, to where alerts are sent.
What kind of endpoint depends on the specific Provider type being used.
For the generic Provider, for example, this is an HTTP/S address.
For other Provider types this could be a project ID or a namespace.</p>
</td>
</tr>
<tr>
<td>
<code>timeout</code><br>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Timeout for sending alerts to the Provider.</p>
</td>
</tr>
<tr>
<td>
<code>proxy</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Proxy the HTTP/S address of the proxy server.
Deprecated: Use ProxySecretRef instead. Will be removed in v1.</p>
</td>
</tr>
<tr>
<td>
<code>proxySecretRef</code><br>
<em>
<a href="https://pkg.go.dev/github.com/fluxcd/pkg/apis/meta#LocalObjectReference">
github.com/fluxcd/pkg/apis/meta.LocalObjectReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ProxySecretRef specifies the Secret containing the proxy configuration
for this Provider. The Secret should contain an &lsquo;address&rsquo; key with the
HTTP/S address of the proxy server. Optional &lsquo;username&rsquo; and &lsquo;password&rsquo;
keys can be provided for proxy authentication.</p>
</td>
</tr>
<tr>
<td>
<code>secretRef</code><br>
<em>
<a href="https://pkg.go.dev/github.com/fluxcd/pkg/apis/meta#LocalObjectReference">
github.com/fluxcd/pkg/apis/meta.LocalObjectReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>SecretRef specifies the Secret containing the authentication
credentials for this Provider.</p>
</td>
</tr>
<tr>
<td>
<code>serviceAccountName</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ServiceAccountName is the name of the service account used to
authenticate with services from cloud providers. An error is thrown if a
static credential is also defined inside the Secret referenced by the
SecretRef.</p>
</td>
</tr>
<tr>
<td>
<code>certSecretRef</code><br>
<em>
<a href="https://pkg.go.dev/github.com/fluxcd/pkg/apis/meta#LocalObjectReference">
github.com/fluxcd/pkg/apis/meta.LocalObjectReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>CertSecretRef specifies the Secret containing TLS certificates
for secure communication.</p>
<p>Supported configurations:
- CA-only: Server authentication (provide ca.crt only)
- mTLS: Mutual authentication (provide ca.crt + tls.crt + tls.key)
- Client-only: Client authentication with system CA (provide tls.crt + tls.key only)</p>
<p>Legacy keys &ldquo;caFile&rdquo;, &ldquo;certFile&rdquo;, &ldquo;keyFile&rdquo; are supported but deprecated. Use &ldquo;ca.crt&rdquo;, &ldquo;tls.crt&rdquo;, &ldquo;tls.key&rdquo; instead.</p>
</td>
</tr>
<tr>
<td>
<code>suspend</code><br>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Suspend tells the controller to suspend subsequent
events handling for this Provider.</p>
</td>
</tr>
<tr>
<td>
<code>commitStatusExpr</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>CommitStatusExpr is a CEL expression that evaluates to a string value
that can be used to generate a custom commit status message for use
with eligible Provider types (github, gitlab, gitea, bitbucketserver,
bitbucket, azuredevops). Supported variables are: event, provider,
and alert.</p>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="notification.toolkit.fluxcd.io/v1beta3.AlertSpec">AlertSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#notification.toolkit.fluxcd.io/v1beta3.Alert">Alert</a>)
</p>
<p>AlertSpec defines an alerting rule for events involving a list of objects.</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>providerRef</code><br>
<em>
<a href="https://pkg.go.dev/github.com/fluxcd/pkg/apis/meta#LocalObjectReference">
github.com/fluxcd/pkg/apis/meta.LocalObjectReference
</a>
</em>
</td>
<td>
<p>ProviderRef specifies which Provider this Alert should use.</p>
</td>
</tr>
<tr>
<td>
<code>eventSeverity</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>EventSeverity specifies how to filter events based on severity.
If set to &lsquo;info&rsquo; no events will be filtered.</p>
</td>
</tr>
<tr>
<td>
<code>eventSources</code><br>
<em>
<a href="https://pkg.go.dev/github.com/fluxcd/notification-controller/api/v1#CrossNamespaceObjectReference">
[]github.com/fluxcd/notification-controller/api/v1.CrossNamespaceObjectReference
</a>
</em>
</td>
<td>
<p>EventSources specifies how to filter events based
on the involved object kind, name and namespace.</p>
</td>
</tr>
<tr>
<td>
<code>inclusionList</code><br>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>InclusionList specifies a list of Golang regular expressions
to be used for including messages.</p>
</td>
</tr>
<tr>
<td>
<code>eventMetadata</code><br>
<em>
map[string]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>EventMetadata is an optional field for adding metadata to events dispatched by the
controller. This can be used for enhancing the context of the event. If a field
would override one already present on the original event as generated by the emitter,
then the override doesn&rsquo;t happen, i.e. the original value is preserved, and an info
log is printed.</p>
</td>
</tr>
<tr>
<td>
<code>exclusionList</code><br>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ExclusionList specifies a list of Golang regular expressions
to be used for excluding messages.</p>
</td>
</tr>
<tr>
<td>
<code>summary</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Summary holds a short description of the impact and affected cluster.
Deprecated: Use EventMetadata instead.</p>
</td>
</tr>
<tr>
<td>
<code>suspend</code><br>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Suspend tells the controller to suspend subsequent
events handling for this Alert.</p>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="notification.toolkit.fluxcd.io/v1beta3.ProviderSpec">ProviderSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#notification.toolkit.fluxcd.io/v1beta3.Provider">Provider</a>)
</p>
<p>ProviderSpec defines the desired state of the Provider.</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>type</code><br>
<em>
string
</em>
</td>
<td>
<p>Type specifies which Provider implementation to use.</p>
</td>
</tr>
<tr>
<td>
<code>interval</code><br>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Interval at which to reconcile the Provider with its Secret references.
Deprecated and not used in v1beta3.</p>
</td>
</tr>
<tr>
<td>
<code>channel</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Channel specifies the destination channel where events should be posted.</p>
</td>
</tr>
<tr>
<td>
<code>username</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Username specifies the name under which events are posted.</p>
</td>
</tr>
<tr>
<td>
<code>address</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Address specifies the endpoint, in a generic sense, to where alerts are sent.
What kind of endpoint depends on the specific Provider type being used.
For the generic Provider, for example, this is an HTTP/S address.
For other Provider types this could be a project ID or a namespace.</p>
</td>
</tr>
<tr>
<td>
<code>timeout</code><br>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Timeout for sending alerts to the Provider.</p>
</td>
</tr>
<tr>
<td>
<code>proxy</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Proxy the HTTP/S address of the proxy server.
Deprecated: Use ProxySecretRef instead. Will be removed in v1.</p>
</td>
</tr>
<tr>
<td>
<code>proxySecretRef</code><br>
<em>
<a href="https://pkg.go.dev/github.com/fluxcd/pkg/apis/meta#LocalObjectReference">
github.com/fluxcd/pkg/apis/meta.LocalObjectReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ProxySecretRef specifies the Secret containing the proxy configuration
for this Provider. The Secret should contain an &lsquo;address&rsquo; key with the
HTTP/S address of the proxy server. Optional &lsquo;username&rsquo; and &lsquo;password&rsquo;
keys can be provided for proxy authentication.</p>
</td>
</tr>
<tr>
<td>
<code>secretRef</code><br>
<em>
<a href="https://pkg.go.dev/github.com/fluxcd/pkg/apis/meta#LocalObjectReference">
github.com/fluxcd/pkg/apis/meta.LocalObjectReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>SecretRef specifies the Secret containing the authentication
credentials for this Provider.</p>
</td>
</tr>
<tr>
<td>
<code>serviceAccountName</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ServiceAccountName is the name of the service account used to
authenticate with services from cloud providers. An error is thrown if a
static credential is also defined inside the Secret referenced by the
SecretRef.</p>
</td>
</tr>
<tr>
<td>
<code>certSecretRef</code><br>
<em>
<a href="https://pkg.go.dev/github.com/fluxcd/pkg/apis/meta#LocalObjectReference">
github.com/fluxcd/pkg/apis/meta.LocalObjectReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>CertSecretRef specifies the Secret containing TLS certificates
for secure communication.</p>
<p>Supported configurations:
- CA-only: Server authentication (provide ca.crt only)
- mTLS: Mutual authentication (provide ca.crt + tls.crt + tls.key)
- Client-only: Client authentication with system CA (provide tls.crt + tls.key only)</p>
<p>Legacy keys &ldquo;caFile&rdquo;, &ldquo;certFile&rdquo;, &ldquo;keyFile&rdquo; are supported but deprecated. Use &ldquo;ca.crt&rdquo;, &ldquo;tls.crt&rdquo;, &ldquo;tls.key&rdquo; instead.</p>
</td>
</tr>
<tr>
<td>
<code>suspend</code><br>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Suspend tells the controller to suspend subsequent
events handling for this Provider.</p>
</td>
</tr>
<tr>
<td>
<code>commitStatusExpr</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>CommitStatusExpr is a CEL expression that evaluates to a string value
that can be used to generate a custom commit status message for use
with eligible Provider types (github, gitlab, gitea, bitbucketserver,
bitbucket, azuredevops). Supported variables are: event, provider,
and alert.</p>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<div class="admonition note">
<p class="last">This page was automatically generated with <code>gen-crd-api-reference-docs</code></p>
</div>
