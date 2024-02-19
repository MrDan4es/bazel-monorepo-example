load("@rules_pkg//:pkg.bzl", _pkg_tar = "pkg_tar")
load("@rules_oci//oci:defs.bzl", _oci_image = "oci_image", _oci_tarball = "oci_tarball")
load("@rules_multirun//:defs.bzl", _command = "command", _multirun = "multirun")

NONROOT = 65532

DEFAULT_BASE = "@distroless_base_debian12"

def container_bundle(name, registry, images):
    [
        _oci_tarball(
            name = "tarball_{}".format(v.split("/")[-1]),
            image = "{}:oci".format(v),
            repo_tags = ["{}/{}".format(registry, k)],
        )
        for k, v in images.items()
    ]

    [
        _command(
            name = "cmd_{}".format(v.split("/")[-1]),
            arguments = [],
            command = ":tarball_{}".format(v.split("/")[-1]),
        )
        for _, v in images.items()
    ]

    _multirun(
        name = name,
        commands = [
            "cmd_{}".format(v.split("/")[-1])
            for _, v in images.items()
        ],
        jobs = 0,
    )

def go_image(name, base = None, embed = [], pure = "on", entrypoint = "/app", **kwargs):
    if not base:
        base = DEFAULT_BASE

    volumes = kwargs.setdefault("volumes", [])
    tars = kwargs.setdefault("tars", [])

    if volumes:
        volumes_tar_name = name + "-volumes-tar"
        _pkg_tar(
            name = volumes_tar_name,
            owner = "%d.%d" % (NONROOT, NONROOT),
            mode = "0700",
            empty_dirs = volumes,
        )
        tars += [":" + volumes_tar_name]

    _pkg_tar(
        name = "layer",
        srcs = embed,
    )

    tars += [":layer"]

    _oci_image(
        name = "oci",
        base = base,
        entrypoint = [entrypoint],
        tars = tars,
        visibility = ["//visibility:public"],
    )
