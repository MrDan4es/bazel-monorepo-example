load("//build/oci:bundle.bzl", _container_bundle = "container_bundle")
load("//build/oci:push_all.bzl", _container_push_all = "container_push_all")
load("//build/oci:image.bzl", _go_image = "go_image")

container_bundle = _container_bundle
container_push_all = _container_push_all
go_image = _go_image
