<p align="center">
    <img alt="Securepoint" title="Securepoint" src="assets/logo.svg" width="100px" height="100px">
</p>

# Securepoint Cloud Shield Definitions

A collection of Cloud Shield definitions, including services with their domains.

## Overview

This repository contains JSON definitions used by Cloud Shield. Services are defined with their associated domains.

## Repository Structure

```
cloud-shield-definitions/
├── README.md
├── schemas/
│   ├── services.schema.json  # JSON schema defining the service structure
├── main.go             # Script to compile all definitions into generated JSON files
├── ids.json            # Stable numeric IDs grouped by service filename
├── services/           # Directory containing all service definitions
│   ├── service1.json
│   ├── service2.json
│   └── ...
```

## Compiled Definitions

All individual definition files are automatically compiled into generated JSON catalogs. Service definitions are compiled into `services.json`.

### Downloading the Compiled File

The compiled files are available in two ways:

1. **GitHub Releases**: Automatically published with each push to master branch
2. **GitHub Actions Artifacts**: Available for all CI runs (including pull requests)

### Local Compilation

To compile the definitions locally:

```bash
go run main.go
```

This will generate:

- `services.json` with the compiled service catalog
- `ids.json` with the stable numeric ID registry in the form `{ "services": {...}`

Each compiled service includes an `id` field. Existing IDs are preserved, and only newly added definition files receive the next free number.

## Adding New Services

To add a new service:

1. Create a new JSON file in the `services` directory
2. Follow the structure defined in `schemas/services.schema.json`
3. Include all required fields and adhere to the specified formats
4. Use a descriptive filename that clearly identifies the service
5. Validate your JSON against the schema before submitting

> [!WARNING]
> **DO NOT change existing file names once they are merged into the repository.**
>
> File names serve as unique identifiers within Cloud Shield.
> Pull requests that change file names will be rejected.

## Contributing

We welcome contributions to this project!
If you have suggestions for new services or improvements, please open an issue or submit a pull request.

### Reporting Issues

- Use GitHub Issues to report problems with existing services
- Include detailed information about the issue

## License

MIT License
