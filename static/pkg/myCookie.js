const myCookie = {
    versionDate : '2025-01-27',
    
    on: navigator.cookieEnabled,
    set: (key, value) => {
        key = String(key)
        value = String(value)
        
        value = value.replace("\n\n", ' ')
        value = value.replace("\n", ' ')
        
        const endAt = new Date()
        endAt.setFullYear(endAt.getFullYear() + 7)
        
        document.cookie = encodeURIComponent(key) + '=' + encodeURIComponent(value) + '; ' + 
        (window.location.protocol === 'https:' ? 'SameSite=None; Secure; ' : '') + 
        'path=/; domain=' + window.location.host + '; expires=' + endAt.toUTCString()
    },
    get: (key) =>  {
        key = String(key)
        key = encodeURIComponent(key)
        let matches = document.cookie.match(new RegExp(
        "(?:^|; )" + key.replace(/([\.$?*|{}\(\)\[\]\\\/\+^])/g, '\\$1') + "=([^;]*)"
        ))
        return matches ? decodeURIComponent(matches[1]) : undefined
    },
    check: (key) => {
        key = String(key)
        key = encodeURIComponent(key)
        return document.cookie.split(';').some((cookie) => {
            return cookie.trim().startsWith(key + '=')
        })
    },
    remove: (key) => {
        key = String(key)
        key = encodeURIComponent(key)
        
        const endAt = new Date(0)
        document.cookie = key + '=; SameSite=None; Secure; path=/; domain=' + 
        window.location.host + '; expires=' + endAt.toUTCString()
    },
    clear: () => {
        let arrayCookie = document.cookie.split(';')
        let name = ''
        while (name = arrayCookie.pop())
            myCookie.remove((name.split('=')[0]).trim())
    },
    show: () => {
        console.log(document.cookie.split(';').map(function(value) { 
            return value.split('=').map(function(value1) {
                return value1.trim()
            }).join(' = ')
        }).join("\n\n"))
    }
}

