/*
need
<div id="notify"></div>
*/

const notify = {
    on:false,
    object: null,
    message: text => {
        notify.add(text, 'message', 4000)
    },
    err: text => {
        notify.add(text, 'err', 4000)
    },
    
    add:(text='HW', type='message', time=4000) => {
        if (!notify.on) {
            notify.start()
            
            setTimeout(() => {
                notify.add(text, type, time)
            }, 2000)
        }
        
        console.log(`notify/${type}/\n${text}`)
        const div = document.createElement('div')
        div.className = `item item-${type}`
        div.innerHTML = text;
        notify.object.insertAdjacentElement('beforeend', div)
        
        setTimeout(() => {
            notify.delete(div)
        }, time)
    },
    start:() => {
        if (notify.on) return;
        notify.on = true
        
        notify.object = document.createElement('div')
        notify.object.id = 'notify'
        document.body.insertAdjacentElement('beforeend', notify.object)
        
        const link = document.createElement('link');
        link.rel = 'stylesheet';
        link.type = 'text/css';
        link.href = './pkg/notify.css?' + Date.now()
        document.head.appendChild(link);
    },
    delete: (element) => {
        element.parentNode.removeChild(element)
    },
    tests:() => {
        notify.start()
        notify.message('Привет! Это обычное уведомление.')
        notify.err('Ошибка! Что-то пошло не так.')
    },
}