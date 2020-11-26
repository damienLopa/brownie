### Brownie

## Build micro-services really fast

# Available default templates

- Flask
- Django

# Usage

brownie -h

# Use custom template

```bash
brownie --template="https://github.com/your_template"
```

# Example of template config


```bash
.brownie.json
{
    "package_installer": [
        "poetry",
        "pipfile"
    ],
    "use_docker":[
        "yes",
        "no"
    ],
    "use_aws_s3": [
        "yes",
        "no"
    ]
}
```
