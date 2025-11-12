const utils = {
    versionDate: '2024-10-31',
    setText: (text = 'Hello', id = 'app') => {
        const element = document.getElementById(id)
        if (element) element.innerHTML = text
    },
    toInt: text => {
        text = String(text)
        text = text.replace(/\D+/g, '')
        
        n = parseInt(text)
        if (Number.isNaN(n)) 
            n = 0
            
        return n
    },
    focusById: id => {
        const element = document.getElementById(id)
        if (element) element.focus()
    },

    parseJson: (jsonString = '{}') => {
        try {
            return JSON.parse(jsonString)
        } catch (error) {
            console.log(['err parseJson', error, jsonString])
        }
        return {}
    }, // JSON.stringify(object)

    
    windowOnLoadAdd: event => {
       if (window.onload)
          window.onload = window.onload + event
       else
          window.onload = event
    },
    handleResponse: (callback = console.log, data, arguments1) => {
        if (data === undefined) 
            callback()
        else if (arguments1 === undefined) 
            callback(data)
        else 
            callback(data, arguments1)
    },
    readInput: id => {
        const element = document.getElementById(id)
        return element ? element.value : '';
    },
    writeInput: (id, text) => {
        const element = document.getElementById(id)
        if (element) element.value = text
    },
    delete: (element) => {
        element.parentNode.removeChild(element)
    },
    
    clearName: name => {
        name = String(name).replace(/[^a-zA-Z0-9+=\-()*?:%;№!.,{}[]@!#%^&*_ ]/g, '')
        name = name.replace(/[<>]/g, '')
        if (name.length > 64) name = name.substring(0, 64)
        name = name.trim()
        return name;
    },
    date: () => {
        const now = new Date()
    
        const year = now.getFullYear()
        const month = String(now.getMonth() + 1).padStart(2, '0')
        const day = String(now.getDate()).padStart(2, '0')
    
        const hours = String(now.getHours()).padStart(2, '0')
        const minutes = String(now.getMinutes()).padStart(2, '0')
        const seconds = String(now.getSeconds()).padStart(2, '0')
    
        return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
    },
    
    // преобразует данные формы в обьект
    formData: form => {
        const formData = new FormData(form);
        let formObject = {};
        
        formData.forEach((value, key) => {
            formObject[key] = value
        })
        
        return formObject
    },
    
    urlParams: () => (
        window
        .location
        .search
        .replace('?','')
        .split('&')
        .reduce(
            (result, line) => {
                let tmp = line.split('=')
                result[decodeURIComponent(tmp[0])] = decodeURIComponent(tmp[1])
                return result;
            },
            {}
        )
    ),
    formatPhone: phone => {
        phone = String(phone).trim()
    
        if (phone === '+') return '+'
    
        const havePlus = phone.indexOf("+") === -1 ? false : true
        phone = phone.replace(/\D/g, '')
    
        const length = phone.length;
        if (length === 0) return ''
    
        // Форматирование для российских номеров (11 цифр)
        if (length === 11) {
            if (phone[0] === '8') phone = '7' + phone.slice(1)
            return `+${phone[0]} (${phone.slice(1, 4)}) ${phone.slice(4, 7)}-${phone.slice(7, 9)}-${phone.slice(9, 11)}`
        }
    
        // Форматирование для других номеров
        let result = ''
        for (let i = 0; i < length; i++) {
            if (i > 0 && i % 3 === 1) result += '-'
            result += phone[i]
        }
    
        return (havePlus ? '+' : '') + result.slice(0, 32)
    },
    isValidEmail: email => {
        const emailRegex = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$/;
        return emailRegex.test(email);
    },
    escapeHtml: text => {
        const map = {
            '&': '&amp;',
            '<': '&lt;',
            '>': '&gt;',
            '"': '&quot;',
            "'": '&#039;'
        }
        return String(text).replace(/[&<>"']/g, m => map[m])
    }
}