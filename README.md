<p align="center">
    <img alt="Securepoint" title="Securepoint" src="assets/logo.svg" width="100px" height="100px">
</p>


# Securepoint DNS Services

A collection of services with the corresponding domains used in our product "Cloud Shield".

## Overview

This repository contains JSON definitions for various services and their associated domains that are used by Cloud Shield. Each service is defined in a separate JSON file that follows a standardized schema to ensure consistency and reliability.

## Repository Structure

```
dns-services/
├── README.md
├── schema.json         # JSON schema defining the service structure
├── main.go             # Script to compile all services into a single JSON file
├── service_ids.json    # Stable numeric IDs per service filename
└── services/           # Directory containing all service definitions
    ├── service1.json
    ├── service2.json
    └── ...
```

## Compiled Services

All individual service files are automatically compiled into a single `services.json` file through the CI/CD pipeline. This compiled file contains all services in a single JSON object where each key is the filename (without `.json` extension).

### Downloading the Compiled File

The compiled `services.json` file is available in two ways:

1. **GitHub Releases**: Automatically published with each push to master branch
2. **GitHub Actions Artifacts**: Available for all CI runs (including pull requests)

### Local Compilation

To compile the services locally:

```bash
go run main.go
```

This will generate:

- `services.json` with the compiled service catalog
- `service_ids.json` with the stable numeric ID registry used by the compiler

Each compiled service now includes an `id` field. Existing IDs are preserved, and only newly added service files receive the next free number.

## Adding New Services

To add a new service:

1. Create a new JSON file in the `services` directory
2. Follow the structure defined in `schema.json`
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
