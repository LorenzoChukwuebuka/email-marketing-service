const fs = require('fs')

//input and validation
const filePath = './README.md'
if (!fs.existsSync(filePath) || !fs.statSync(filePath).isFile()) {
    console.error('File does not exist or is not readable.')
    process.exit(1)
}

// 2. Open File
const buffer = fs.readFileSync(filePath)
const chunkSize = 16
let offset = 0

//3 process file chunks

for (let i = 0; i < buffer.length; i += chunkSize) {
    const chunk = buffer.slice(i, i + chunkSize) // Get chunk
    const hex = []
    const ascii = []

    // 4. Process Each Byte in the Chunk
    chunk.forEach(byte => {
        hex.push(byte.toString(16).padStart(2, '0')) // Convert to hex
        ascii.push(byte >= 32 && byte <= 126 ? String.fromCharCode(byte) : '.') // ASCII or '.'
    })

    // 5. Format Output
    const hexStr = hex.join(' ')
    const asciiStr = ascii.join('')
    console.log(
        `${offset.toString(16).padStart(8, '0')}: ${hexStr.padEnd(
            47
        )} ${asciiStr}`
    )

    offset += chunkSize // Increment offset
}


 
