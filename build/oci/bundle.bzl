load("@rules_oci//oci:defs.bzl", _oci_tarball = "oci_tarball")
load("@rules_multirun//:defs.bzl", _command = "command", _multirun = "multirun")

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
