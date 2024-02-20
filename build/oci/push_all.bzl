load("@rules_oci//oci:defs.bzl", _oci_push = "oci_push", _oci_tarball = "oci_tarball")
load("@rules_multirun//:defs.bzl", _command = "command", _multirun = "multirun")

def container_push_all(name, images, registry):
    [
        _oci_push(
            name = "push_{}".format(v.split("/")[-1]),
            image = "{}:oci".format(v),
            repository = "{}/{}".format(registry, k.split(":")[0]),
            remote_tags = ["latest"],
        )
        for k, v in images.items()
    ]

    [
        _command(
            name = "cmd_push_{}".format(v.split("/")[-1]),
            command = ":push_{}".format(v.split("/")[-1]),
            arguments = [],
        )
        for _, v in images.items()
    ]

    _multirun(
        name = name,
        commands = [
            "cmd_push_{}".format(v.split("/")[-1])
            for _, v in images.items()
        ],
        jobs = 0,
    )
