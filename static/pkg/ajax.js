const ajax = {
    versionDate: '2025-07-03',
    handleResponse: (callback = console.log, data, arg) => {
        if (data === undefined) 
            callback()
        else if (arg === undefined) 
            callback(data)
        else 
            callback(data, arg)
    },
    json: (url, data = {}, callback = console.log, arg) => {
        if (!url) 
            throw new Error('URL are required.')
    
        fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        })
        .then(response => {
            if (!response.ok) 
                throw new Error(response.status + ' : ' + response.statusText)
            
            return response.json()   
        })
        .then(responseJson => {
            ajax.handleResponse(callback, responseJson, arg)
        })
        .catch(error => {
            console.error(error)
        })
    },
    jsonTest: (url, data = {}, callback = console.log, arg) => {
        console.log({
            'function':'ajax.jsonTest',
            'url':url,
            'data':data,
            'callback':callback,
            'arg':arg
        })
        if (!url)
            throw new Error('URL are required.')

        fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        })
        .then(response => {
            if (!response.ok)
                throw new Error(response.status + ' : ' + response.statusText)
            return response.text()
        })
        .then(responseText => {
            console.log(responseText)
            try {
                const responseJson = JSON.parse(responseText)
                ajax.handleResponse(callback, responseJson, arg)
            } catch (error) {
                console.error('Ошибка парсинга JSON:', error)
            }
        })
        .catch(error => {
            console.error(error)
        });
    },
    getJson: (url, callback = console.log, arg) => {
        if (!url) 
            throw new Error('URL are required.')
    
        fetch(url)
        .then(response => {
            if (!response.ok) 
                throw new Error(response.status + ' : ' + response.statusText)

            return response.json()   
        })
        .then(responseJson => {
            ajax.handleResponse(callback, responseJson, arg)
        })
        .catch(error => {
            console.error(error)
        });
    },
    text: (url, callback = console.log, arg) => {
        if (!url) 
            throw new Error('URL are required.')
            
        fetch(url)
        .then(response => response.text())
        .then(text => {
            ajax.handleResponse(callback, text, arg)
        })
        .catch(error => {
            console.error(error)
        })
    },
    sendForm: (url, form, callback = console.log, arg) => {
        if (!url)
            throw new Error('URL are required.')

        let formData = new FormData(form);

        fetch(url, {
            method: 'POST',
            body: formData
        })
        .then(response => {
            if (!response.ok)
                throw new Error(response.status + ' : ' + response.statusText)
            return response.json()
        })
        .then(responseJson => {
            ajax.handleResponse(callback, responseJson, arg)
        })
        .catch(error => {
            console.error(error)
        })
    },
    
}