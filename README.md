<!--

add resume

https://github.com/golang-standards/project-layout

 -->

# [g.harel.page](https://g.harel.page)

&nbsp;

![diagram](https://user-images.githubusercontent.com/9319710/48919035-6c882b80-ee5e-11e8-87b2-4d15f4f8f866.png)

&nbsp;

This website is generated every five minutes by a [Cloud Function](https://cloud.google.com/functions/) which uploads the result to a public [Cloud Storage Bucket](https://cloud.google.com/storage/) whose contents are served by [Cloudflare](https://www.cloudflare.com/). Each build fetches the contents of [.config](./.config) to build a [GraphQL](https://graphql.org/) query for the [GitHub API](https://developer.github.com/v4/). This data is then used to execute the website [templates](./templates) and build the static output files.

This process means the website contents, like user icon and repo descriptions, will always be up-to-date, without impacting response time and scalability. It also makes it easy, fast and version-controlled to update the list of projects and contributions shown on the website.

## Development

```
$ go run ./development/build.go
```

This command will output `index.html` in the project root and uses the local `.config` to fetch for data.

After the initial build, the `./templates` directory is watched for changes. Manual rebuilds can be triggered by typing `.\n`.

If a `.env` file exists in the project root, it will automatically be loaded.

The `GRAPHQL_TOKEN` environment variable must be defined to access the GitHub API.

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
* _[Default cloud function service account](https://cloud.google.com/functions/docs/concepts/iam#runtime_service_account) must have `Storage Object Admin` Role._

## Planned Improvements

* Continuously `terraform apply` repository contents using [Cloud Build](https://cloud.google.com/cloud-build/)
* Make styling responsive
* Better secret management using [KMS](https://cloud.google.com/kms/)

## License

[MIT](https://github.com/g-harel/website/blob/master/LICENSE)
