<!--

nodemon -e go,html,css,js,svg -x "dotenv go run ./development/build.go"

add comments
add CV link + hosted
add analytics?
add mobile banner

https://github.com/golang-standards/project-layout

 -->

# [g.harel.page](https://g.harel.page)

&nbsp;

![diagram](https://user-images.githubusercontent.com/9319710/48919035-6c882b80-ee5e-11e8-87b2-4d15f4f8f866.png)

&nbsp;

The website is generated every five minutes by a [Cloud Function](https://cloud.google.com/functions/) which uploads the result to a public [Cloud Storage Bucket](https://cloud.google.com/storage/) whose contents are served by [Cloudflare](https://www.cloudflare.com/). Each build fetches contents of [.config](./.config) to build a [GraphQL](https://graphql.org/) query for the [GitHub API](https://developer.github.com/v4/). This data is then used to execute the website [templates](./templates) and build the static output files.

This process means the website contents, like user icon and repo descriptions, will always be up-to-date, without impacting response time and scalability. It also makes it easy, fast and version-controlled to update the list of projects and contributions shown on the website.

## Development

```
$ GRAPHQL_TOKEN="{github_api_token}" go run ./development/build.go
```

This will use local `.config` and output `index.html` in the project root.

## Deployment

```
$ make build
```

```
$ cd ./terraform
$ terraform apply
```

_deployment is not required for changes to `.config`_

_As of this writing (Nov 2018)_

* _Unreleased version of [terraform-provider-google](https://github.com/terraform-providers/terraform-provider-google) is required to allow non-node Cloud Function runtime. See [here](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins) for installation instructions after building the provider from source._
* _Project owner and terraform service account must be whitelisted for the go Cloud Function alpha._
* _Cloud Scheduler resource must be managed manually because it is not yet supported in [terraform-provider-google](https://github.com/terraform-providers/terraform-provider-google)._

## Planned Improvements

* Continuously `terraform apply` repository contents using [Cloud Build](https://cloud.google.com/cloud-build/)
* Move source templates into storage bucket to avoid redeploying function on template change
* Make styling responsive
* Better secret management (possibly using [Vault](https://www.vaultproject.io/))

## License

[MIT](https://github.com/g-harel/website/blob/master/LICENSE)
