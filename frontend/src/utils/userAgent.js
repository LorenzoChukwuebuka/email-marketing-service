const parseUserAgent = userAgent => {
    // Example parsing logic
    const info = {}
  
    if (userAgent.includes('Chrome')) {
      info.name = 'Google Chrome'
    } else if (userAgent.includes('Firefox')) {
      info.name = 'Mozilla Firefox'
    } else if (userAgent.includes('Safari') && !userAgent.includes('Chrome')) {
      info.name = 'Apple Safari'
    } else if (userAgent.includes('Edge')) {
      info.name = 'Microsoft Edge'
    } else if (userAgent.includes('MSIE') || userAgent.includes('Trident/')) {
      info.name = 'Internet Explorer'
    } else {
      info.name = 'Unknown Browser'
    }
  
    // Extract version
    const versionMatch = userAgent.match(
      /(Chrome|Firefox|Safari|Edge|MSIE|rv:|Opera|OPR)\/([\d.]+)/
    )
    if (versionMatch && versionMatch.length >= 3) {
      info.version = versionMatch[2]
    } else {
      info.version = 'Unknown'
    }
  
    return info
  }
  
  export default parseUserAgent