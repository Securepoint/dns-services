const fs = require("node:fs")
const path = require("node:path")

const serviceDir = path.join(__dirname, "services")
const outputFile = path.join(__dirname, "services.json")

const services = {}

try {
    const files = fs.readdirSync(serviceDir)
    const jsonFiles = files.filter((file) => file.endsWith(".json"))

    console.log(`found ${jsonFiles.length} service files`)

    for (const file of jsonFiles) {
        const filePath = path.join(serviceDir, file)
        const fileContent = fs.readFileSync(filePath, "utf8")

        try {
            const serviceData = JSON.parse(fileContent)
            const serviceName = path.basename(file, ".json")

            services[serviceName] = serviceData

            console.log(`compiled: ${serviceName}`)
        } catch (parseError) {
            console.error(`error parsing ${file}:`, parseError.message)
            process.exit(1)
        }
    }

    fs.writeFileSync(outputFile, JSON.stringify(services, null, 4), "utf8")

    console.log(
        `successfully compiled ${Object.keys(services).length} services to ${path.basename(outputFile)}`
    )
} catch (error) {
    console.error("error compiling services:", error.message)
    process.exit(1)
}
