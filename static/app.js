/*



*/
const options = {
    version: '9 2025-09-04',
    apiUrl: '/api/v2/',
}
var storage = {
    idb: null,
    tgBot:''
}

const views = {
    header: (data={}) => (`
<nav id="mainMenu" class="mainMenu">
    <div class="mainMenuLogo">
        <img src="./images/logo-white.svg" alt="logo">
        <h2>Trade-In</h2>
    </div>
    ${data.clear ? '' : `
    <h4>Продавцам</h4>
    <button onmousedown="pages.clients()"><img src="./images/ico/person-vcard.svg" alt="ico">Клиенты</button>
    <button onmousedown="pages.refunds()"><img src="./images/ico/arrow-left-right.svg" alt="ico">Возвраты</button>
    
    <h4>Внутренние</h4>
    <button onmousedown="pages.shopGroups()"><img src="./images/ico/award.svg" alt="ico">Сети</button>
    <button onmousedown="pages.shops()"><img src="./images/ico/shop.svg" alt="ico">Магазины</button>
    <button onmousedown="pages.shipments()"><img src="./images/ico/boxes.svg" alt="ico">Поставки</button>

    <h4></h4>
    <button onmousedown="pages.users()"><img src="./images/ico/people.svg" alt="ico">Пользователи</button>
    <button onmousedown="pages.profile()"><img src="./images/ico/person.svg" alt="ico">Профиль</button>
    <button onmousedown="pages.help()"><img src="./images/ico/question-octagon.svg" alt="ico">Помощь</button>
    `}
</nav>
    `),
    /*
        <button onmousedown="pages.products()"><img src="./images/ico/shop.svg" alt="ico">Товары</button>
    <button onmousedown="pages.productsCategories()"><img src="./images/boxes.svg" alt="ico">Категории</button>
    <button onmousedown="pages.productsColors()"><img src="./images/palette.svg" alt="ico">Цвета товаров</button>
    <button onmousedown="pages.questions()"><img src="./images/question-octagon.svg" alt="ico">Вопросы</button>
    <button onmousedown="pages.texts()"><img src="./images/body-text.svg" alt="ico">Текста</button>
    */
    
    listSelectCity: data => (
        data.map(item => (`
    <button class="btn-selectCityItem" type="button" ontouchstart="app.selectCity(${item.id})" onmousedown="app.selectCity(${item.id})">${item.name} (${item.region})</button>
        `))
        .join('')
    ),
    search: data => (`
<form 
    onsubmit="app.search(this);return false" 
    class="search-form"
>
<input 
    class="search-input"
    id="searchInput" 
    type="text" 
    value=""
>
<button 
    class="button search-button"
></button>
</form>
    `),
    listKnows: data => (`
<ul class="know-list">
${storage.list.map(know => (`
    <li>
        <button class="know-list-name">${know.name}</button>
    </li>
`)).join('')}
</ul>    
    `),
    addKnowButton: data => (`
<button 
    class="button addknow-button"
    onmousedown="pages.addKnow()"
>Добавить</button>    
    `),
    addKnowForm: data => (`
    при изменениях сверху выезжает сохранить
<textarea></textarea>
<button class="button">Настройки</button>
<!--button class="button">Тип текст</button-->
    `),
}
let pages = {
    help: data => {
        utils.setText(` 
        ${views.header(data)}
        <div class="content" id="content">
            <div class="title block">
                <h1>Помощь</h1>
                <div class="text">
                    <p>
                    <b>Как работает программа?</b><br>
                    Покупатель приносит свое старое устройство для участия в программе по traid-in в партнёрский магазин. <span>На обмен принимаются все pod-системы - бренд и год выпуска не важны (работающую/не работающую, главное что бы она не была разбита, вздута).</span><br><br>
                    
                    <span>Покупателю необходимо купить 1 картридж examp2 5.</span><br> После покупки одного картриджа, продавец регистрирует клиента на <span>seller.example.ru.</span><br>
                    Вносит все необходимые обязательные поля, удостоверившись, что на момент покупки покупателю исполнилось 18 лет.<br><br>
                    
                    После успешной регистрации покупателю выдается тушка  examp2 5.<br><br>
                    
                    ТРЕЙД-ИН ВОЗМОЖНО ПРОИЗВЕСТИ ТОЛЬКО 1 РАЗ.<br>
                    ОДНА ТУШКА  examp2 5 — В ОДНИ РУКИ.<br>
                    ПО ОДНОМУ И ТОМУ ЖЕ НОМЕРУ ТЕЛЕФОНА ПРОГРАММА НЕ ДАСТ ПРОВЕСТИ ЕЩЕ ОДИН ОБМЕН!!<br><br>
                    
                    Старое устройство - утилизируйте сами.<br><br>
                    
                    Перечень устройств указан при регистрации нового клиента. В разделе — устройство обмена. Если устройства нет в перечне — выбираете «другое».<br><br>
                    </p>
                    
                    <p><b>Контакты</b><br> 
                        
                        
                        
                        <br><br>
                        <b>Где брать код устройства?</b><br><img class="width100" src="./images/instruction.jpg" alt="instruction">
                    </p>
                </div>
                <div class="buttons">
                    <button class="btn btn-mark" onmousedown="pages.formAdd()">Администратор</button>
                    <button class="btn" onmousedown="pages.formAdd()">Менеджер</button>
                    <button class="btn" onmousedown="pages.formAdd()">Продавец</button>
                    <button class="btn" onmousedown="pages.formAdd()">Новый пользователь</button>
                </div>
            </div>
            
        </div>
        `)
    },

    login: data => {
      utils.setText(`
<div class="loginBox">
    <h1 class="logo_text">Trade-In</h1>
    <form class="form1 block" id="formLogin" onsubmit="app.loginPhone(this);return false">
        <label class="name">Телефон</label>
        <div class="value"><input name="phone" oninput="this.value = utils.formatPhone(this.value)" type="text" value="" maxlength="64"></div>
        <!--<br>
        label class="name required">Пароль</label>
        <div class="value"><input name="password" type="password" value="" maxlength="64"></div-->
        <div class="buttons">
            <button class="btn btn-mark">Войти</button>
        </div>
        <div class="buttons">
            <button class="btn" type="button" onmousedown="pages.loginRules()">Регистрация</button>
            <button class="btn" type="button" onmousedown="app.loginTest()">Тест</button>
        </div>
    </form>
</div>
      `);
    },
    loginRules: data => {
        utils.setText(`
<div class="loginBox">
    <h1 class="logo_text">Trade-In</h1>
    <form class="form1 block" id="formLogin" onsubmit="pages.login();app.registrationTG();return false">
        <label class="name">Правила</label>
        <div class="value">Попадёте в телеграм бот, там нажмите старт, после отправить номер. Тогда, при входе, можно будет вводить свой номер</div>
        <div class="form-check">
          <input type="checkbox" id="policy_check" required="">
          <label for="policy_check" class="form-check-label">Я даю свое <a href="/policy/privacy_consent2.pdf" target="_blank">согласие на обработку персональных данных</a> и соглашаюсь с <a href="/policy/policy2.pdf" target="_blank">условиями политики конфиденциальности</a></label>
        </div>
        <!--<br>
        label class="name required">Пароль</label>
        <div class="value"><input name="password" type="password" value="" maxlength="64"></div-->
        <div class="buttons">
            <button class="btn btn-mark">Согласен</button>
            <button class="btn" type="button" onmousedown="pages.login()">Назад</button>
        </div>
    </form>
</div>
      `);
    },

    loginCode: data => {
        utils.setText(`
<div class="loginBox">
    <h1 class="logo_text">Trade-In</h1>
    <form class="form1 block" id="formLogin" onsubmit="app.loginCode(this);return false">
        <label class="name">Введите код</label>
        <div class="value"><input name="code" oninput="this.value = utils.toInt(this.value)" type="text" value="" maxlength="6"></div>
        
        <div class="value">Код Вам пришёл в телеграм бот, привязанный к номеру </div>
        <!--<br>
        label class="name required">Пароль</label>
        <div class="value"><input name="password" type="password" value="" maxlength="64"></div-->
        <div class="buttons">
            <button class="btn btn-mark">Вход</button>
            <button class="btn" type="button" onmousedown="pages.login()">Назад</button>
        </div>
    </form>
</div>
      `);
    },
    
    clients: data => {
        utils.setText(` 
${views.header(data)}
<div class="content" id="content">
    <div class="title block">
        <h1>Клиенты</h1>
        <div class="text">
            <p>Список клиентов, которые вы добавили</p>
            <p>После добавления, отсутствует возможность удалить и изменить. Если допустили ошибку, то обратитесь к администратору, контакты указаны на странице помощи</p>
        </div>
        <div class="buttons">
            <button class="btn btn-mark" onmousedown="pages.client()"><img src="./images/ico/plus-circle.svg" alt="ico">Добавить клиента</button>
        </div>
    </div>
    <div class="list block" style="grid-template-columns: 1fr 3fr 2fr 1fr 2fr;">
        <div class="line">
            <div class="grid-item grid-header">ID</div>
            <div class="grid-item grid-header">ФИО Клиента</div>
            <div class="grid-item grid-header">Телефон</div>
            <div class="grid-item grid-header">Код продукта</div>
            <div class="grid-item grid-header">Устройство</div>
        </div>
    
        
    </div>
</div>
        `)
        
    },
    client: data => {
        if (!data) data = {};
        
        // Если передан ID клиента — подгружаем данные из хранилища
        if (
            data.id && 
            storage.data.Client.mapIdIndex[data.id] !== undefined && 
            storage.data.Client.list[storage.data.Client.mapIdIndex[data.id]]
        ) {
            data = storage.data.Client.list[storage.data.Client.mapIdIndex[data.id]];
        }
    
        // Опции для устройства обмена (из вашего HTML)
        const productOptions = [
            { value: '63', text: 'VAPORESSO' },
            { value: '64', text: 'example' },
            { value: '65', text: 'VOOPOO' },
            { value: '66', text: 'GEEKVAPE' },
            { value: '67', text: 'OXVA' },
            { value: '68', text: 'SMOANT' },
            { value: '69', text: 'RINCOE' },
            { value: '70', text: 'LOST VAPE' },
            { value: '0', text: 'Другое' }
        ];
    

    
        utils.setText(`
            ${views.header(data)}
            <div class="content" id="content">
                <div class="title block">
                    <h1>Добавить клиента</h1>
                    <div class="text">
                        <p>Поля с * обязательны для заполнения.</p>
                        <p>Город - выбирается автоматически текущий город магазина</p>
                    </div>
                </div>
    
                <form class="form1 block" onsubmit="return false">
                    <label class="name required">Магазин (в котором вы сейчас)</label>
                    <div class="value">
                        <select name="shop" id="shop" required>
                            <option value="" disabled selected>Выбрать</option>
                            <option value="5">Адрес 1</option>
                            <option value="20">Адрес id20</option>
                            <option value="21">Адрес</option>
                            <option value="30">Адрес</option>
                            <option value="40">Адрес длинный длинный длинный длинный длинный длинный длинный 555</option>
                        </select>
                    </div>
                    
                    <!--label class="name">id: ${data.id || 0}</label>
                    <div class="value"><input name="id" type="hidden" value="${data.id || 0}"></div-->
    
                    <label class="name required">Фамилия</label>
                    <div class="value">
                        <input
                            type="text"
                            name="family_name"
                            value="${data.family_name || ''}"
                            placeholder="Введите фамилию"
                            maxlength="128"
                            required
                        >
                    </div>
    
                    <label class="name required">Имя</label>
                    <div class="value">
                        <input
                            type="text"
                            name="name"
                            value="${data.name || ''}"
                            placeholder="Введите имя"
                            maxlength="128"
                            required
                        >
                    </div>
                    
                    <label class="name required">Отчество</label>
                    <div class="value">
                        <input
                            type="text"
                            name="middle_name"
                            value="${data.middle_name || ''}"
                            placeholder="Введите отчество"
                            maxlength="128"
                            required
                        >
                    </div>
    
                    <label class="name required">Телефон</label>
                    <div class="value">
                        <input
                            type="tel"
                            name="phone"
                            value="${data.phone || ''}"
                            onfocus="app.formatPhoneNumber(this)"
                            oninput="app.formatPhoneNumber(this)"
                            placeholder="+7 (XXX) XXX-XX-XX"
                            required
                        >
                    </div>
    
                    <label class="name">Email</label>
                    <div class="value">
                        <input
                            type="email"
                            name="email"
                            value="${data.email || ''}"
                            placeholder="Введите E-mail"
                            maxlength="256"
                        >
                    </div>
    

    
                    <label class="name required">Дата рождения</label>
                    <div class="value">
                        <input
                            type="text"
                            name="birth_date"
                            value="${data.birth_date || ''}"
                            placeholder="дд/мм/гггг"
                            required
                        >
                    </div>
    
                    <label class="name required">Пол</label>
                    <div class="value">
                        <div class="gender">
                            <label class="gender-radio">
                                <input type="radio" id="male" name="gender_id" value="1" ${data.gender_id == 1 ? 'checked' : ''} required> М
                            </label>
                            <label class="gender-radio">
                                <input type="radio" id="female" name="gender_id" value="2" ${data.gender_id == 2 ? 'checked' : ''}> Ж
                            </label>
                        </div>
                    </div>
    
                    <label class="name required">Устройство обмена</label>
                    <div class="value">
                        <select name="old_product_id" id="old_product_id" required>
                            <option value="" disabled selected>Выбрать</option>
                            ${productOptions.map(opt => 
                                `<option value="${opt.value}">${opt.text}</option>`
                            ).join('')}
                        </select>
                                              
                        
                        
                        <input
                            
                            type="text"
                            name="old_product_name"
                            id="old_product_name" 
                            value="${data.old_product_name || ''}" 
                            style="display: none;" 
                            placeholder="Введите название устройства"
                        >
                        
                        
                        
                        
                    </div>
    
                    <label class="name required">Модель устройства</label>
                    <div class="value">
                        <input
                            type="text"
                            name="old_product_model"
                            value="${data.old_product_model || ''}"
                            placeholder=""
                            required
                        >
                    </div>
    
                    <label class="name required">Код продукта</label>
                    <div class="value">
                        <input
                            type="text"
                            name="product_code"
                            value="${data.product_code || ''}"
                            placeholder="BR00000"
                            maxlength="7"
                            oninput="app.formatDeviceInput(this)"
                            required
                        >
                    </div>
    
                    <!-- Картриджи -->
                    <label class="name required">Количество купленных картриджей MINIKAN 5 (если нет, ставим 0)</label>
                    <div class="value">
                        <input
                            type="text"
                            name="quantity_cartridge06"
                            value="${data.quantity_cartridge06 || ''}"
                            oninput="this.value = this.value == '' ? '' :utils.toInt(this.value)"
                            placeholder="0,6 Ом"
                            required
                        >
                        <input
                            type="text"
                            name="quantity_cartridge08"
                            value="${data.quantity_cartridge08 || ''}"
                            oninput="this.value = this.value == '' ? '' :utils.toInt(this.value)"
                            placeholder="0,8 Ом"
                            required
                        >
                        <input
                            type="text"
                            name="quantity_cartridge1"
                            value="${data.quantity_cartridge1 || ''}"
                            oninput="this.value = this.value == '' ? '' :utils.toInt(this.value)"
                            placeholder="1 Ом"
                            required
                        >
                    </div>
                    <div class="buttons">
                        <button class="btn btn-mark">Добавить</button>
                    </div>
                </form>
            </div>
        `);

    },
    
    refunds: data => {
        utils.setText(` 
${views.header(data)}
<div class="content" id="content">
    <div class="title block">
        <h1>Возвраты / Обмен</h1>
        <div class="text">
            <p>Список возвратов, которые вы добавили.</p>
            <p>Когда клиент приходит с проблемой, что выданное через trade-in устройство сломанное, то необходимо забрать старое и выдать ему новое, и оформить возврат.</p>
            <p>Такой возврат/обмен можно делать только 1 раз для каждого клиента.</p>
            <p>Список причин допустимые для возврата: !!!</p>
            <p>Список причин НЕ допустимые для возврата: !!!</p>
        </div>
        <div class="buttons">
            <button class="btn btn-mark" onmousedown="pages.refund()"><img src="./images/ico/plus-circle.svg" alt="ico">Добавить возврат</button>
        </div>
    </div>
    <div class="list block" style="grid-template-columns: 1fr 2fr 4fr 1fr;">
        <div class="line">
            <div class="grid-item grid-header">ID</div>
            <div class="grid-item grid-header">Клиент</div>
            <div class="grid-item grid-header">Причина</div>
            <div class="grid-item grid-header">Код нового продукта</div>
        </div>
    
        
    </div>
</div>
        `)
        
    },
    refund: data => {
        if (!data) data = {};
        
        utils.setText(`
            ${views.header(data)}
            <div class="content" id="content">
                <div class="title block">
                    <h1>Добавить возврат / обмен</h1>
                    <div class="text">
                        <p>Поля с * обязательны для заполнения.</p>
                        <p>Первично необходимо проверить код возвращаемого продукта, и после удачных проверок, вторая форма будет доступна</p>
                    </div>
                </div>
    
                <form class="form1 block" onsubmit="return false">
                    <label class="name required">Код возвращаемого продукта</label>
                    <div class="value">
                        <input
                            type="text"
                            name="refund_product_code"
                            value="${data.refund_product_code || ''}"
                            placeholder="BR00000"
                            maxlength="7"
                            oninput="app.formatDeviceInput(this)"
                            required
                        >
                    </div>
                    <div class="buttons">
                        <button class="btn btn-mark" onmousedown="pages.refund()">Проверить!</button>
                    </div>
                </form>
                <form class="form1 block lock" onsubmit="return false">
                    <label class="name required">Магазин (в котором вы сейчас)</label>
                    <div class="value">
                        <select name="shop" id="shop" required>
                            <option value="" disabled selected>Выбрать</option>
                            <option value="5">Адрес 1</option>
                            <option value="20">Адрес id20</option>
                            <option value="21">Адрес</option>
                            <option value="30">Адрес</option>
                            <option value="40">Адрес длинный длинный длинный длинный длинный длинный длинный 555</option>
                        </select>
                    </div>

                    <label class="name required">Причина</label>
                    <div class="value"> 
                        <textarea name="reason" required>${data.reason || ''}</textarea>
                    </div>
                    
                    <label class="name required">Код нового продукта</label>
                    <div class="value">
                        <input
                            type="text"
                            name="product_code"
                            value="${data.product_code || ''}"
                            placeholder="BR00000"
                            maxlength="7"
                            oninput="app.formatDeviceInput(this)"
                            required
                        >
                    </div>
                    <div class="buttons">
                        <button class="btn btn-mark">Добавить</button>
                    </div>
                </form>
            </div>
        `);

    },
    
    
    shopGroups: data => {
        utils.setText(` 
${views.header(data)}
<div class="content" id="content">
    <div class="title block">
        <h1>Сети</h1>
        <div class="text">
        </div>
        <div class="buttons">
            <button class="btn btn-mark" onmousedown="pages.shopGroup()"><img src="./images/ico/plus-circle.svg" alt="ico">Добавить сеть</button>
        </div>
    </div>
    <div class="list block" style="grid-template-columns: 1fr 5fr 1fr;">
        <div class="line">
            <div class="grid-item grid-header">ID</div>
            <div class="grid-item grid-header">Название</div>
            <div class="grid-item grid-header">Действия</div>
        </div>
    </div>
</div>
        `)
        
    },
    shopGroup:data => {
        if (!data) data = {};
        
        utils.setText(`
            ${views.header(data)}
            <div class="content" id="content">
                <div class="title block">
                    <h1>Добавить сеть</h1>
                    <div class="text">
                        <p>Поля с * обязательны для заполнения.</p>
                    </div>
                </div>
    
                <form class="form1 block" onsubmit="return false">
                    <label class="name required">Название</label>
                    <div class="value">
                        <input
                            type="text"
                            name="name"
                            value="${data.name || ''}"
                            placeholder="Название сети"
                            maxlength="64"
                            required
                        >
                    </div>
                    <div class="buttons">
                        <button class="btn btn-mark">Добавить</button>
                    </div>
                </form>
            </div>
        `);
    },
    
    shops: data => {
        utils.setText(` 
${views.header(data)}
<div class="content" id="content">
    <div class="title block">
        <h1>Магазины</h1>
        <div class="text">
            <p>Список магазинов вашей сети</p>
        </div>
        <div class="buttons">
            <button class="btn btn-mark" onmousedown="pages.shop()"><img src="./images/ico/plus-circle.svg" alt="ico">Добавить магазин</button>
            <button class="btn" onmousedown="pages.shops()">Обновить список</button>
        </div>
    </div>
    <div class="list block" style="grid-template-columns: 1fr 5fr 1fr 1fr;">
        <div class="line">
            <div class="grid-item grid-header">ID</div>
            <div class="grid-item grid-header">Адрес</div>
            <div class="grid-item grid-header">Наличие устройств</div>
            <div class="grid-item grid-header">Действия</div>
        </div>
    </div>
</div>
        `)
        
    },
    shop:data => {
        if (!data) data = {};
        
        utils.setText(`
            ${views.header(data)}
            <div class="content" id="content">
                <div class="title block">
                    <h1>Добавить магазин</h1>
                    <div class="text">
                        <p>Поля с * обязательны для заполнения.</p>
                        <p>Магазин добавляется только вашей сети. У админа есть возможность выбрать сеть</p>
                    </div>
                </div>
    
                <form class="form1 block" onsubmit="return false">
                    <label class="name required">Адрес</label>
                    <div class="value">
                        <input
                            type="text"
                            name="address"
                            value="${data.address || ''}"
                            placeholder="Адрес"
                            maxlength="256"
                            required
                        >
                    </div>
                    <div class="buttons">
                        <button class="btn btn-mark">Добавить</button>
                    </div>
                </form>
            </div>
        `);
    },
    
    shipments: data => {
        utils.setText(` 
${views.header(data)}
<div class="content" id="content">
    <div class="title block">
        <h1>Поставки</h1>
        <div class="text">
            <p>Список клиентов, которые вы добавили</p>
        </div>
        <div class="buttons">
            <button class="btn btn-mark" onmousedown="pages.shipment()"><img src="./images/ico/plus-circle.svg" alt="ico">Добавить поставку</button>
            <button class="btn" onmousedown="">Обновить список</button>
        </div>
    </div>
    <div class="list block" style="grid-template-columns: 1fr 2fr 1fr 1fr 1fr;">
        <div class="line">
            <div class="grid-item grid-header">ID</div>
            <!--div class="grid-item grid-header">Сеть</div-->
            <div class="grid-item grid-header">Магазин</div>
            <div class="grid-item grid-header">Количество</div>
            <div class="grid-item grid-header">Статус</div>
            <div class="grid-item grid-header">Действия</div>
        </div>
    
        
    </div>
</div>
        `)
        
    },
    shipment:data => {
        if (!data) data = {};
        
        utils.setText(`
            ${views.header(data)}
            <div class="content" id="content">
                <div class="title block">
                    <h1>Добавить поставку</h1>
                    <div class="text">
                        <p>Поля с * обязательны для заполнения.</p>
                    </div>
                </div>
    
                <form class="form1 block" onsubmit="return false">
                    <label class="name required">Сеть</label>
                    <div class="value">
                        <select name="shop" id="shop" required>
                            <option value="" disabled selected>Выбрать</option>
                            <option value="5">Сеть 1</option>
                            <option value="20">Сеть id20</option>
                            <option value="21">Сеть</option>
                            <option value="30">Сеть</option>
                            <option value="40">Сеть 555</option>
                        </select>
                    </div>
                    <label class="name required">Магазин</label>
                    <div class="value">
                        <select name="shop" id="shop" required>
                            <option value="" disabled selected>Выбрать</option>
                            <option value="5">Адрес 1</option>
                            <option value="20">Адрес id20</option>
                            <option value="21">Адрес</option>
                            <option value="30">Адрес</option>
                            <option value="40">Адрес длинный длинный длинный длинный длинный длинный длинный 555</option>
                        </select>
                    </div>
                    
                    <label class="name required">Количество</label>
                    <div class="value">
                        <input
                            type="text"
                            name="name"
                            value="${data.name || ''}"
                            placeholder="Название сети"
                            maxlength="64"
                            required
                        >
                    </div>
                    <label class="name required">Статус  !!! Уточнить сптсок статусов</label>
                    <div class="value">
                        <select name="shop" id="shop" required>
                            <option value="" disabled selected>Выбрать</option>
                            <option value="5">Подготовка к отправке</option>
                            <option value="20">Отправили</option>
                            <option value="21">Получен</option>
                            <option value="30">Отменён</option>
                        </select>
                    </div>
                    
                    <label class="name required">Описание</label>
                    <div class="value"> 
                        <textarea name="reason" required>${data.description || ''}</textarea>
                    </div>
                    
                    <div class="buttons">
                        <button class="btn btn-mark">Добавить</button>
                    </div>
                </form>
            </div>
        `);
    },
    
    users: data => {
        utils.setText(` 
${views.header(data)}
<div class="content" id="content">
    <div class="title block">
        <h1>Пользователи</h1>
        <div class="text">
            <p>Список пользователей данной системы (не клиенты)</p>
            <p>Доступно только админам и менеджерам</p>
            <p>Нет возможности добавить, так как регистрируются они сами и попадают в список Свободны</p>
            <p>Из списка свободных им определяют роль</p>
            <p>Роль админа выдаётся только через обращения по контакту, которые в помощи</p>
            <p>Менеджер может быть назначен за определённой сетью и выдавать роли только этой сети, так же видеть кто в его сети</p>
        </div>
        <div class="buttons">
            <button class="btn" onmousedown="">Свободные</button>
            <button class="btn" onmousedown="">Продавцы</button>
            <button class="btn" onmousedown="">Менеджеры</button>
            <button class="btn" onmousedown="">Админы</button>
        </div>
    </div>
    <div class="list block" style="grid-template-columns: 1fr 3fr 2fr 1fr 1fr;">
        <div class="line">
            <div class="grid-item grid-header">ID</div>
            <div class="grid-item grid-header">ФИО</div>
            <div class="grid-item grid-header">Телефон</div>
            <div class="grid-item grid-header">Роль</div>
            <div class="grid-item grid-header">Действия</div>
        </div>
    
        
    </div>
</div>
        `)
        
    },
    user:data => {},
    
    profile:data => {
        if (!data) data = {};
        
        utils.setText(`
            ${views.header(data)}
            <div class="content" id="content">
                <div class="title block">
                    <h1>Профиль</h1>
                    <div class="text">
                        <p>Менять только через админа. Контакты в помощи</p>
                    </div>
                </div>
    
                <form class="form1 block" onsubmit="return false">
                    <label class="name"><b>ФИО:</b> Фам имя отч</label>
                    <label class="name"><b>Роль:</b> Продавец</label>
                    <label class="name"><b>Клиентов:</b> 20</label>
                    <label class="name"><b>Возвратов:</b> 20</label>
                    
                    <label class="name">Ваш магазин по умолчанию</label>
                    <div class="value">
                        <select name="shop" id="shop" required>
                            <option value="" disabled selected>Выбрать</option>
                            <option value="5">Адрес 1</option>
                            <option value="20">Адрес id20</option>
                            <option value="21">Адрес</option>
                            <option value="30">Адрес</option>
                            <option value="40">Адрес длинный 555</option>
                        </select>
                    </div>
                </form>
            </div>
        `);
    },
    
    
    
    registration: (data) => {
        if (!data) data = {}
        utils.setText(`
            ${views.header({clear:true})}
            <div class="content" id="content">
                <div class="title block">
                    <h1>Регистрация</h1>
                    <div class="text">
                        <p>Только для тех, кому больше 18 лет.</p>
                        <p>После регистрации вам должны будут дать должность, после чего доступы придут на почту</p>
                    </div>
                </div>
    
                <form class="form1 block" onsubmit="return false">
                    <label class="name required">Фамилия</label>
                    <div class="value">
                        <input
                            type="text"
                            name="family_name"
                            value="${data.family_name || ''}"
                            placeholder="Введите фамилию"
                            maxlength="128"
                            required
                        >
                    </div>
    
                    <label class="name required">Имя</label>
                    <div class="value">
                        <input
                            type="text"
                            name="name"
                            value="${data.name || ''}"
                            placeholder="Введите имя"
                            maxlength="128"
                            required
                        >
                    </div>
                    
                    <label class="name required">Отчество</label>
                    <div class="value">
                        <input
                            type="text"
                            name="middle_name"
                            value="${data.middle_name || ''}"
                            placeholder="Введите отчество"
                            maxlength="128"
                            required
                        >
                    </div>
    
                    <label class="name required">Телефон</label>
                    <div class="value">
                        <input
                            type="tel"
                            name="phone"
                            value="${data.phone || ''}"
                            onfocus="app.formatPhoneNumber(this)"
                            oninput="app.formatPhoneNumber(this)"
                            placeholder="+7 (XXX) XXX-XX-XX"
                            required
                        >
                    </div>
    
                    <label class="name">Email</label>
                    <div class="value">
                        <input
                            type="email"
                            name="email"
                            value="${data.email || ''}"
                            placeholder="Введите E-mail"
                            maxlength="256"
                        >
                    </div>
    

    
                    <label class="name required">Дата рождения</label>
                    <div class="value">
                        <input
                            type="text"
                            name="birth_date"
                            value="${data.birth_date || ''}"
                            placeholder="дд/мм/гггг"
                            required
                        >
                    </div>
    
                    <label class="name required">Пол</label>
                    <div class="value">
                        <div class="gender">
                            <label class="gender-radio">
                                <input type="radio" id="male" name="gender_id" value="1" ${data.gender_id == 1 ? 'checked' : ''} required> М
                            </label>
                            <label class="gender-radio">
                                <input type="radio" id="female" name="gender_id" value="2" ${data.gender_id == 2 ? 'checked' : ''}> Ж
                            </label>
                        </div>
                    </div>
    
                    <label class="name required">Согласие</label>
                    <div class="value">
                        <div class="form-check">
                          <input type="checkbox" id="age_check" required>
                          <label for="age_check" class="form-check-label">Мне есть 18 лет</label>
                        </div>
                        <div class="form-check">
                          <input type="checkbox" id="policy_check" required>
                          <label for="policy_check" class="form-check-label">Я даю свое <a href="/policy/privacy_consent2.pdf" target="_blank">согласие на обработку персональных данных</a> и соглашаюсь с <a href="/policy/policy2.pdf" target="_blank">условиями политики конфиденциальности</a></label>
                        </div>
                    </div>

                    <div class="buttons">
                        <button class="btn btn-mark">Зарегистрироваться</button>
                    </div>
                    
                    <div class="buttons">
                        <p>Уже зарегистрированы?</p>
                        <button class="btn" onmousedown="pages.login()">Войти</button>
                    </div>
                </form>
            </div>
        `)
        window.scrollTo(0,0)
    },
    
    
    
    
    
    // ---
    main: (data) => {
      if (!storage.user || !storage.user.id) {
        return pages.login();
      }
      if (!data || !data.stop)
        app.ajax("list", {}, (answer) => {
          if (answer.err) return app.showPopup(answer.err, true);
  
          storage.list = answer;
          pages.main({ stop: 1 });
        });

        // <p class="descr">При регистрации каждого клиента продавец получает 100 баллов.<br>(баллы выводятся через <a href="https://partners.example.ru/">partners.example.ru</a>)</p>
        utils.setText(` 
          ${app.addHeader()}
          <div class="main page">
              <h1 class="title">Зарегистрируйте клиента!</h1>
              <div class="tabs">
                  <button class="btn" onmousedown="pages.formAdd()">Добавить клиента</button>
                  <button class="btn btn-transp" onmousedown="pages.searchUserByBrCode()">
                    Оформить возврат
                  </button>
              </div>
              <div class="table">
                  <div class="table-header">
                      <p class="table-col">ID</p>
                      <p class="table-col">Имя клиента</p>
                      <p class="table-col">Телефон</p>
                      <p class="table-col">Код продукта</p>
                      <p class="table-col">Устройство</p>
                  </div>
                  ${storage.list ? storage.list.map(item => `
                      <div class="user-list-item">
                          <div class="table-row">
                            <p class="table-label">ID:</p>
                            <p class="table-value">${item.id}</p>
                          </div>
                          <div class="table-row">
                            <p class="table-label">Имя клиента:</p>
                            <p class="table-value">${item.family_name} ${item.name}</p>
                          </div>
                          <div class="table-row">
                            <p class="table-label">Телефон:</p>
                            <p class="table-value">${item.phone}</p>
                          </div>
                          <div class="table-row">
                            <p class="table-label">Код продукта:</p>
                            <p class="table-value">${item.product_code}</p>
                          </div>
                          <div class="table-row">
                            <p class="table-label">Устройство:</p>
                            <p class="table-value ${item.product_name && item.product_name.split(' ').length > 1 ? item.product_name.split(' ')[1].toLowerCase() : ''}">${item.product_name && item.product_name.split(' ').length > 1 ? item.product_name.split(' ')[1] : ''}</p>
                          </div>
                      </div>`).join("") : ""}
              </div>
              <div id="popup">
                  <span id="close-btn" onmousedown="app.closePopup()"><svg id="closeSvg" width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                      <path d="M7.25367 18C6.93501 18 6.61635 17.8826 6.36478 17.6312C5.87841 17.145 5.87841 16.3402 6.36478 15.854L15.8574 6.36465C16.3438 5.87845 17.1488 5.87845 17.6352 6.36465C18.1216 6.85086 18.1216 7.65561 17.6352 8.14181L8.14256 17.6312C7.90776 17.8826 7.57233 18 7.25367 18Z" fill="#A2E9C1"/>
                      <path d="M16.7463 18C16.4277 18 16.109 17.8826 15.8574 17.6312L6.36478 8.14181C5.87841 7.65561 5.87841 6.85086 6.36478 6.36465C6.85115 5.87845 7.65618 5.87845 8.14256 6.36465L17.6352 15.854C18.1216 16.3402 18.1216 17.145 17.6352 17.6312C17.3836 17.8826 17.065 18 16.7463 18Z" fill="#A2E9C1"/>
                  </svg></span>
                  <p id="popup-message"></p>
              </div>
          </div>
      `);
      window.scrollTo(0,0);
    },
    
    account: (data) => {
      utils.setText(`
          ${app.addHeader()}
          <div class="page">
          <div class="container1400">
          <button class="back-btn" onmousedown="pages.main()">Назад</button>
            ${['1'].includes(storage.user.group_id) ? `
                <button class="back-btn" onmousedown="app.request.all()" >Заявки</button>
            ` : ''}
          <div class="profile">
            <div class="title">Профиль</div>
            <div class="seller-info">
              <p><span class="info-title">ФИО: </span>${storage.user.family_name} ${storage.user.name} ${storage.user.middle_name}</p>
              <p><span class="info-title">Дата рождения: </span>${storage.user.birth_date}</p>
              <p><span class="info-title">Телефон: </span>${storage.user.phone}</p>
              <p><span class="info-title">Email: </span>${storage.user.email}</p>
              <p><span class="info-title">Город: </span>${storage.cities.find(city => city.id === storage.user.city_id).name}</p>
              <p><span class="info-title">Адрес магазина: </span>${storage.user.address}</p>
              <p><span class="info-title">Название магазина: </span>${storage.user.shop}</p>
            </div>
            <button class="btn btn-primary" onmousedown="app.logout()">Выйти из аккаунта</button>
          </div>
          </div>
          </div>
          `)
      window.scrollTo(0,0);
    },
    list: (data) => {
      if (!data || !data.stop)
        app.ajax("list", {}, (answer) => {
          if (answer.err) return app.showPopup(answer.err, true);
  
          storage.list = answer;
          pages.list({ stop: 1 });
        });
  
      utils.setText(`
              <div class="page">
                  <div class="title">Список</div>
                  Место для мотивирующей статистики, рейтинга<br><br>
                  <div class="tabs">
                      <button class="tab-button" onmousedown="pages.main()"><i class="bi bi-house-door-fill"></i></button>
                      <button class="tab-button" onmousedown="pages.formAdd()">Добавить</button>
                      <button class="tab-button" onmousedown="pages.account()"><i class="bi bi-person-circle me-2"></i>${
                        storage.user.name
                      }</button>
                  </div>
                  ${
                    storage.list
                      ? storage.list
                          .map(
                            (item) => `
                  <div class="user-list-item">
                      [${item.id}] ${item.family_name} ${item.name} ${item.middle_name}<br>${item.phone}<br>${item.email}<br>created_at ${item.created_at}<br>${item.product_code}
                  </div>`
                          )
                          .join("")
                      : ""
                  }
                  <div id="popup">
                      <span id="close-btn" onmousedown="app.closePopup()"><svg id="closeSvg" width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M7.25367 18C6.93501 18 6.61635 17.8826 6.36478 17.6312C5.87841 17.145 5.87841 16.3402 6.36478 15.854L15.8574 6.36465C16.3438 5.87845 17.1488 5.87845 17.6352 6.36465C18.1216 6.85086 18.1216 7.65561 17.6352 8.14181L8.14256 17.6312C7.90776 17.8826 7.57233 18 7.25367 18Z" fill="#A2E9C1"/>
<path d="M16.7463 18C16.4277 18 16.109 17.8826 15.8574 17.6312L6.36478 8.14181C5.87841 7.65561 5.87841 6.85086 6.36478 6.36465C6.85115 5.87845 7.65618 5.87845 8.14256 6.36465L17.6352 15.854C18.1216 16.3402 18.1216 17.145 17.6352 17.6312C17.3836 17.8826 17.065 18 16.7463 18Z" fill="#A2E9C1"/>
</svg>
</span>
                      <p id="popup-message"></p>
                  </div>
              </div>
          `);
    },
    wait: (data) => {
      utils.setText(`Ждите`);
    },
    formAdd: (data) => {
      utils.setText(`
        ${app.addHeader()}
        <div class="page">
        <div class="flex-box">
        <button class="back-btn" onmousedown="pages.main()">Назад</button>
        <div class="card">
            <h2 class="card-header">Новый клиент</h2>
            <div class="card-body">
              <form onsubmit="${
        storage.user.group_id == "3"
          ? `app.addClient(this);`
          : `app.submitForm(this, 'add');`
      } return false;">
          <div class="form-group">
              <label for="family_name">Фамилия</label>
              <input type="text" name="family_name" id="family_name"  placeholder="Введите фамилию" onblur="app.handleBlur(this)" oninput="app.restrictInput(this)" required>
          </div>

                    <div class="form-group">
                        <label for="name">Имя</label>
                        <input type="text" name="name" id="name"  placeholder="Введите имя"
                            onblur="app.handleBlur(this)" oninput="app.restrictInput(this)" required>
                    </div>

                    <div class="form-group">
                        <label for="middle_name">Отчество</label>
                        <input type="text" name="middle_name" id="middle_name" 
                            placeholder="Введите отчество" onblur="app.handleBlur(this)"
                            oninput="app.restrictInput(this)" required>
                    </div>

                    <div class="form-group">
                        <label for="phone">Телефон</label>
                        <input type="tel" name="phone" id="phone" 
                            onfocus="app.startPhoneNumber(this)" oninput="app.formatPhoneNumber(this)"
                            placeholder="+7 (XXX) XXX-XX-XX" required>
                    </div>

                    <div class="form-group">
                        <label class="mail-label" for="email">Email</label>
                        <input type="email" name="email" id="email"  placeholder="Введите E-mail">
                    </div>

                    <div class="form-group">
                        <label for="city_name">Город</label>
                        <input type="text" name="city_name" id="city_name" 
                            placeholder="Введите город" autocomplete="off"
                           
                            oninput="app.handleCityInput(this)" 
                             
                            required>
                        <input type="hidden" name="city_id" id="city_id">
                        <div id="citySelect" style="display: none;"></div>
                        <div id="city-error" class="error-message" style="display: none;">Введите корректное название
                        </div>
                    </div>

                    ${
                    storage.user.group_id == "2"
                    ? `
                    <div class="form-group">
                        <label for="address">Адрес магазина</label>
                        <input type="text" name="address" id="address" 
                            placeholder="Введите Адрес магазина" oninput="app.handleAddressInput(event)" required>
                        <div id="suggestions" style="display: none;"></div>
                    </div>

                    <div class="form-group">
                        <label for="shop">Название магазина</label>
                        <input type="text" name="shop" id="shop" 
                            placeholder="Введите Название магазина" required>
                    </div>
                    `
                    : ""
                    }

                    ${
                    storage.user.group_id == "3"
                    ? `
                    <input type="hidden" name="action" value="add">
                    <div class="form-group">
                        <label for="birth_date">Дата рождения</label>
                        <input type="text" name="birth_date" id="birth_date" 
                            placeholder="дд/мм/гггг" oninput="app.formatDateInput(this)" onblur="app.validateDateInput(this)" required>
                        <span id="age-error" class="error-message" style="display: none;">Вам должно быть не менее 18
                            лет</span>
                    </div>

                    <div class="form-group">
                        <label class="label-gender">Пол</label>
                        <div class="gender">
                            <label class="gender-radio"><input type="radio" id="male" name="gender_id" value="1" required>М</label>
                            <label class="gender-radio"><input type="radio" id="female" name="gender_id" value="2">Ж</label>
                        </div>
                    </div>
                  
                    

                    <div class="form-group custom-select">
                      <label for="old_product_id">Устройство обмена</label>
                      <div class="selected-option form-control" onmousedown="app.toggleDropdown()">
                        <span id="selected-text">Выберите устройство</span>
                        <span class="arrow">&#9660;</span>
                        <input type="text" name="old_product_id" id="old_product_id" value="" required>
                      </div>
                      <div class="options" id="options">
                        <div class="option" data-value="">Выберите устройство</div>
                        ${storage.products ? storage.products.map(item => `<div class="option" data-value="${item.id}">${item.name}</div>`).join('') : ''}
                        <div class="option" data-value="0">Другое</div>
                      </div>
                      <input  type="text" name="old_product_name" id="old_product_name" value="" style="display: none;" placeholder="Введите название устройства">
                    </div>
                    
                    <div class="form-group">
                        <label for="old_product_model">Модель</label>
                        <input type="text" name="old_product_model" id="old_product_model" 
                            placeholder="" required>
                    </div>



                    <div class="form-group">
                        <label for="product_code">Код продукта</label>
                        <input type="text" name="product_code" id="product_code" 
                            placeholder="BR00000" maxlength="7" oninput="app.formatDeviceInput(this)" required>
                    </div>

                    <!--<div class="form-group">
                        <label for="passport_number">Номер паспорта</label>
                        <input type="text" name="passport_number" id="passport_number" 
                            placeholder="0000000000" maxlength="10" oninput="app.formatPassportNumInput(this)" required>
                    </div>

                    <div class="form-group scan">
                      <label class="btn btn-transp btn-scan">Прикрепить скан паспорта
                        <input type="file" name="passport_photos[]" id="passport_photos" accept="image/*" multiple required onchange="app.updateFileNames()">
                      </label>
                      <small id="file-names"></small>
                    </div>--!>

                    <div class="form-group">
                        <label for="quantity_cartridge06" >Количество купленных Картриджей MINIKAN 5</label>

                        <input type="text" 
                        name="quantity_cartridge06" id="quantity_cartridge06" 
                        oninput="this.value = utils.stringToInt(this.value)" 
                         placeholder="0,6 Ом" required >
                        
                        <input type="text" 
                        name="quantity_cartridge08" id="quantity_cartridge08" 
                        oninput="this.value = utils.stringToInt(this.value)" 
                         placeholder="0,8 Ом" required >
                        
                        <input type="text" 
                        name="quantity_cartridge1" id="quantity_cartridge1" 
                        oninput="this.value = utils.stringToInt(this.value)" 
                         placeholder="1 Ом" required >
                    </div>


                    `
                    : ""
                    }

                    <button type="submit" id="addButton" class="btn form-btn">Добавить</button>
                </form>
                <div id="popup">
                  <span id="close-btn" onmousedown="app.closePopup()"><svg id="closeSvg" width="24" height="24"
                    viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path
                        d="M7.25367 18C6.93501 18 6.61635 17.8826 6.36478 17.6312C5.87841 17.145 5.87841 16.3402 6.36478 15.854L15.8574 6.36465C16.3438 5.87845 17.1488 5.87845 17.6352 6.36465C18.1216 6.85086 18.1216 7.65561 17.6352 8.14181L8.14256 17.6312C7.90776 17.8826 7.57233 18 7.25367 18Z"
                        fill="#A2E9C1" />
                    <path
                        d="M16.7463 18C16.4277 18 16.109 17.8826 15.8574 17.6312L6.36478 8.14181C5.87841 7.65561 5.87841 6.85086 6.36478 6.36465C6.85115 5.87845 7.65618 5.87845 8.14256 6.36465L17.6352 15.854C18.1216 16.3402 18.1216 17.145 17.6352 17.6312C17.3836 17.8826 17.065 18 16.7463 18Z"
                        fill="#A2E9C1" />
                      </svg>
                  </span>
                  <p id="popup-message"></p>
                </div>
            </div>
          </div>
        </div>
    </div>
  
          `)
      window.scrollTo(0,0);
      app.initSelectDevice();
    },
    password: () => {
      utils.setText(`
        <div class="gradient"></div>
        <div class="login">
        <div class="card">
            <h3 class="card-header">Восстановление пароля</h3>
            <div class="card-body">
                <form onsubmit="app.forgotPassword();return false">
                    <div class="form-group">
                        <label for="email" class="form-label">E-mail</label>
                        <input type="email" name="login"  id="email" placeholder="Введите E-mail" required />
                    </div>
                    <button type="submit" class="btn form-btn">Восстановить</button>
                </form>
            </div>
            <div id="popup">
            <span id="close-btn" onmousedown="app.closePopup()"><svg id="closeSvg" width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M7.25367 18C6.93501 18 6.61635 17.8826 6.36478 17.6312C5.87841 17.145 5.87841 16.3402 6.36478 15.854L15.8574 6.36465C16.3438 5.87845 17.1488 5.87845 17.6352 6.36465C18.1216 6.85086 18.1216 7.65561 17.6352 8.14181L8.14256 17.6312C7.90776 17.8826 7.57233 18 7.25367 18Z" fill="#A2E9C1"/>
              <path d="M16.7463 18C16.4277 18 16.109 17.8826 15.8574 17.6312L6.36478 8.14181C5.87841 7.65561 5.87841 6.85086 6.36478 6.36465C6.85115 5.87845 7.65618 5.87845 8.14256 6.36465L17.6352 15.854C18.1216 16.3402 18.1216 17.145 17.6352 17.6312C17.3836 17.8826 17.065 18 16.7463 18Z" fill="#A2E9C1"/></svg>
            </span>
            <p id="popup-message"></p>
        </div>
    </div>
        </div>
      `);
    },
    recoverPassword: () => {
      utils.setText(`
        <div class="gradient"></div>
        <div class="login">
          <div class="card">
            <div class="card-body">
              <div class="pswd-alert">Новые данные для входа отправлены на почту</div>
              <button class="btn form-btn" onmousedown="app.logout()">Войти</button>
            </div>
          </div>
        </div>
      `);
    },
    
    instruction: (data) => {
    utils.setText(`
        ${app.addHeader()}
        <div class="page">
          <div class="container1400">
            <button class="back-btn" onmousedown="pages.main()">Назад</button>
            <div class="instruction">
              <div class="title">Как работает программа?</div>
              <p class="instruction-txt">
                Покупатель приносит свое старое устройство для участия в программе по traid-in в партнёрский магазин. <span>На обмен принимаются все pod-системы.</span><br><br>

                <span>Покупателю необходимо купить 1 картридж examp2 5.</span><br> После покупки двух картриджей, продавец регистрирует клиента на <span>seller.example.ru.</span><br>
                Вносит все необходимые обязательные поля, удостоверившись, что на момент покупки покупателю исполнилось 18 лет.<br><br>
                
                После успешной регистрации покупателю выдается тушка  examp2 5.<br><br>

                ТРЕЙД-ИН ВОЗМОЖНО ПРОИЗВЕСТИ ТОЛЬКО 1 РАЗ.<br>
                ОДНА ТУШКА  examp2 5 — В ОДНИ РУКИ.<br>
                ПО ОДНОМУ И ТОМУ ЖЕ НОМЕРУ ТЕЛЕФОНА ПРОГРАММА НЕ ДАСТ ПРОВЕСТИ ЕЩЕ ОДИН ОБМЕН!!<br><br>

                Старое устройство складываете в коробку.<br><br>
                
                Перечень устройств указан при регистрации нового клиента. В разделе — устройство обмена. Если устройства нет в перечне — выбираете «другое».<br><br>
              </p>
              <img src="/imgs/instruction.png" alt="instruction-img">
              <br>
              <div class="version">${app.version}</div>
            </div>
          </div>
        </div>
        `)
    // Продавец получает за каждого зарегистрированного пользователя — 100 баллов. Которые можно вывести через робобил — <a href="https://partners.example.ru/">partners.example.ru</a><br><br>

    //           Так же баллы можно получить при продаже и регистрации чеков в «робобил» — <a href="https://partners.example.ru/">partners.example.ru</a><br><br>

    //           <span>ЭС example examp2&nbsp;5 и ЭС example examp2&nbsp;5&nbsp;PRO&nbsp;—&nbsp;150&nbsp;баллов.<br>
    //           Бестабачная кальянная смесь example&nbsp;—&nbsp;15&nbsp;баллов.<br>

    //           Табак жевательный MONSTER CHEWER / HAPPMAN / ANGRY CHEW / SKALA — 30 баллов.<br>
    //           Одноразовая ЭС example VINI&nbsp;—&nbsp;75&nbsp;баллов.</span>
      window.scrollTo(0,0);
    },
    searchUserByBrCode: (data) => {
        // Если не продавец и не админ
        if (
            storage.user.group_id != '3' && storage.user.group_id != '1'
        ) return;

        utils.setText(`
            ${app.addHeader()}
            <div class="page">
                <div class="flex-box">
                    <button class="back-btn" onmousedown="pages.main()">Назад</button>
          
                    <div class="card">
                        <h2 class="card-header">Поиск устройства</h2>
                        <div class="card-body">
                            <form onsubmit="app.searchUserByBrCode(this); return false;">
                                <input type="hidden" name="action" value="searchUserByBRCode">

                                <div class="form-group">
                                    <label for="br_code">Серийный номер миникана, который хотят вернуть</label>
                                    <input type="text" name="br_code" id="br_code" 
                                        placeholder="BR00000" maxlength="7" oninput="app.formatDeviceInput(this)" required>
                                </div>

                                <button type="submit" id="addButton" class="btn form-btn">Поиск</button>
                            </form>

                            <div id="popup">
                                <span id="close-btn" onmousedown="app.closePopup()"><svg id="closeSvg" width="24" height="24"
                                    viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                    <path
                                        d="M7.25367 18C6.93501 18 6.61635 17.8826 6.36478 17.6312C5.87841 17.145 5.87841 16.3402 6.36478 15.854L15.8574 6.36465C16.3438 5.87845 17.1488 5.87845 17.6352 6.36465C18.1216 6.85086 18.1216 7.65561 17.6352 8.14181L8.14256 17.6312C7.90776 17.8826 7.57233 18 7.25367 18Z"
                                        fill="#A2E9C1" />
                                    <path
                                        d="M16.7463 18C16.4277 18 16.109 17.8826 15.8574 17.6312L6.36478 8.14181C5.87841 7.65561 5.87841 6.85086 6.36478 6.36465C6.85115 5.87845 7.65618 5.87845 8.14256 6.36465L17.6352 15.854C18.1216 16.3402 18.1216 17.145 17.6352 17.6312C17.3836 17.8826 17.065 18 16.7463 18Z"
                                        fill="#A2E9C1" />
                                    </svg>
                                </span>
                                <p id="popup-message"></p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        `);
    },
    formDefect: (data) => {
        // Если не продавец или админ
        if (
            storage.user.group_id != "3" && storage.user.group_id != "1"
        ) return;

        utils.setText(`
              ${app.addHeader()}
              <div class="page">
                  <div class="flex-box">
                    <button class="back-btn" onmousedown="pages.main()">Назад</button>
    
                    <div class="card">
                      <h2 class="card-header">Оформление брака</h2>
                      <p>Удостовертись в данных по клиенту</p>
                      <div class="card-body">
                        <form onsubmit="app.editClientDefect(this); return false;">
                          <input type="hidden" name="action" value="editDefect">
                          <input type="hidden" name="user_id" value="${data.id}">
                          
                          <div class="form-group">
                              <label for="family_name">Фамилия</label>
                              <input type="text" name="family_name" id="family_name"  value="${data.family_name}" readonly>
                          </div>
    
                          <div class="form-group">
                            <label for="name">Имя</label>
                            <input type="text" name="name" id="name"  placeholder="Введите имя" value="${data.name}" readonly>
                          </div>
    
                          <div class="form-group">
                            <label for="middle_name">Отчество</label>
                            <input type="text" name="middle_name" id="middle_name"  value="${data.middle_name}" readonly>
                          </div>
    
                          <div class="form-group">
                            <label for="phone">Телефон</label>
                            <input type="tel" name="phone" id="phone"  value="${data.phone}" readonly>
                          </div>
    
                          <div class="form-group">
                            <label for="reason_refund">Причиница дефекта</label>
                            <input type="text" name="reason_refund" id="reason_refund"  maxlength="255" required>
                          </div>
    
                          <div class="form-group">
                            <label for="product_code_refund">Серийный номер нового миникана, который отдаём на замену старого</label>
                            <input type="text" name="product_code_refund" id="product_code_refund" 
                              placeholder="BR00000" maxlength="7" oninput="app.formatDeviceInput(this)" required>
                          </div>
    
                          <button type="submit" id="addButton" class="btn form-btn">Сохранить</button>
                        </form>
    
                        <div id="popup">
                          <span id="close-btn" onmousedown="app.closePopup()"><svg id="closeSvg" width="24" height="24"
                            viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                            <path
                                d="M7.25367 18C6.93501 18 6.61635 17.8826 6.36478 17.6312C5.87841 17.145 5.87841 16.3402 6.36478 15.854L15.8574 6.36465C16.3438 5.87845 17.1488 5.87845 17.6352 6.36465C18.1216 6.85086 18.1216 7.65561 17.6352 8.14181L8.14256 17.6312C7.90776 17.8826 7.57233 18 7.25367 18Z"
                                fill="#A2E9C1" />
                            <path
                                d="M16.7463 18C16.4277 18 16.109 17.8826 15.8574 17.6312L6.36478 8.14181C5.87841 7.65561 5.87841 6.85086 6.36478 6.36465C6.85115 5.87845 7.65618 5.87845 8.14256 6.36465L17.6352 15.854C18.1216 16.3402 18.1216 17.145 17.6352 17.6312C17.3836 17.8826 17.065 18 16.7463 18Z"
                                fill="#A2E9C1" />
                              </svg>
                          </span>
                          <p id="popup-message"></p>
                        </div>
                      </div>
                    </div>
                  </div>
              </div>
        `);
    },
    requests: (data) => {
        // Если не продавец или админ
        if (
          storage.user.group_id != "1"
        ) return;
    
        let items = '<h2>На данный момент нерешенных заявок нет</h2>';
        if (data.length > 0) {
          items = data.map((el, i) => {
            return `
              <div class="card">
                <div class="card-title">
                  <h2>#${el.id}</h2>
                </div>
                <div class="card-body">
                  <ul>
                    <li>
                        <strong>Продавец:</strong>
                        <ul style="padding-left: 20px;">
                          <li>
                            <strong>Идентификатор:</strong> ${el.seller_id}
                          </li>
                          <li>
                            <strong>Имя:</strong> ${el.seller_username}
                          </li>
                          <li>
                            <strong>Телефон:</strong> ${el.seller_phone}
                          </li>
                        </ul>
                    </li>
                    <li>
                        <strong>Клиент:</strong>
                        <ul style="padding-left: 20px;">
                          <li><strong>Идентификатор:</strong> ${el.client_id}</li>
                          <li><strong>Имя:</strong> ${el.client_username}</li>
                          <li><strong>Телефон:</strong> ${el.client_phone}</li>
                        </ul>
                    </li>
                    <li><strong>Код продукта:</strong> ${el.product_code}</li>
                    <li><strong>Причина дефекта:</strong> ${el.content}</li>
                    <li><strong>Дата создания:</strong> ${el.created_at}</li>
                  </ul>
                </div>
                <div class="card-footer request-btn">
                  <button class="btn" onmousedown="app.request.apply(${el.id})">Одобрить</button>
                  <button class="btn btn-transp" onmousedown="app.request.denine(${el.id})">Отменить</button>
                </div>
              </div>
          `;
          }).join('');
        }
    
    
        utils.setText(`
            ${app.addHeader()}
            <div class="page">
              <div class="flex-box">
                <button class="back-btn" onmousedown="pages.main()">Назад</button>
                <div>
                  <div class="title">Заявки</div>
                  ${items}
                </div>
              </div>
            </div>
        `);
    },
}

const app = {
    version: options.version,
    main: () => {
        
        pages.login()
        notify.start()
        app.updateTgBot()
/*      let hashArray = location.hash.split('/')
      if (hashArray[0] == '#createNewPassword') {
        const urlWithoutHash = window.location.href.split('#')[0]
        window.history.replaceState({}, document.title, urlWithoutHash)
        
        myCookie.set("userID", Number(hashArray[1]))
        myCookie.set("userKEY", hashArray[2])
        
        app.ajax('createNewPassword') 
        pages.recoverPassword()
        return;
      }
      
      app.checkAuth(() => {
        pages.main();
      })
      app.updateCities()
      app.updateProducts()*/
    },

    login: (form) => {
        const data = utils.formData(form)

        if (data.phone.length != 18 ) {
            notify.err('Введите корректный номер телефона')
            return
        }

        pages.loginCode();


        phone = utils.toInt(data.phone)
        app.ajax('sendAuthCode', {phone:phone}, answer => {
            console.log(answer)
        })

    },
    loginTest: () => {
        pages.help()
        notify.message('Тестовый режим')
    },

    updateTgBot: () => {
        app.ajax('bot', {}, answer => {
            storage.tgBot = answer.bot
        })
    },
    ajax: (action, data, callback, arg) => {
      if (!data) data = {};
      data.action = action;
  
      ajax.json(`${options.apiUrl}${action}`, data, callback, arg);
    },
    checkAuth: (callback) => {
      if (!myCookie.check("userID") || !myCookie.check("userKEY")) {
        callback();
        return;
      }
  
      // Получение пользователя
      app.ajax("checkAuth", null, (answer) => {
        if (answer.auth && answer.auth == 1) {
          answer.fullName = `${answer.name} ${answer.family_name}`;
          storage.user = answer;
        }
        callback();
      });
    },
    logout: () => {
      app.ajax("logout", null, () => {
        myCookie.remove("userID");
        myCookie.remove("userKEY");
        storage.user = null;
        pages.main();
      });
    },
    forgotPassword: () => {
      const email = document.getElementById('email').value;
      app.ajax('forgotPassword', {email: email}, answer => {
        if (answer.err) return app.showPopup('Введите корректный email', true);
        app.showPopup(`Инструкции для восстановления пароля отправлены на почту: ${email}`, false);
        setTimeout(() => (app.logout()), 2000);

      });
    },
    submitForm: (form, action, callback, arg) => {
      if (app.receiver[action]) callback = app.receiver[action];
      app.ajax(action, utils.formData(form), callback, arg);
    },
    receiver: {
      login: (answer) => {
        if (answer.err) return app.showPopup(answer.err, true);

  
        if (
          answer.login &&
          answer.login == 1 &&
          answer.userID &&
          answer.userKEY
        ) {
          myCookie.set("userID", answer.userID);
          myCookie.set("userKEY", answer.userKEY);
          app.main();
        }
      },
      add: (answer) => {
        const addButton = document.getElementById('addButton');
    
        if (answer.err) {
            app.showPopup(answer.err, true);
            addButton.textContent = 'Добавить';
            addButton.disabled = false;
            return;
        }
    
        if (!answer.userID) {
            addButton.textContent = 'Добавить';
            addButton.disabled = false;
            return;
        }
    
        app.showPopup(
            `Пользователь userID[${answer.userID}] успешно добавлен!`,
            false
        );
    
        setTimeout(() => {
            pages.main();
            addButton.textContent = 'Добавить';
            addButton.disabled = false;
        }, 1500);
      },
      registration: (answer) => {
        if (answer.err) return app.showPopup(answer.err, true);
  
        if (!answer.userID) return;
        app.showPopup(
          `Регистрация успешна! 
              userID ${answer.userID}
              Данные для входа ${
                answer.sendMail == 1 ? "" : "НЕ "
              }отправлены на ваш email ${answer.email}`,
          false
        );
        setTimeout(() => {
          pages.login();
        }, 1500);
      },
      editClientDefect: (answer) => {
        const addButton = document.getElementById('addButton');

        if (!answer.ok) {
            app.showPopup(answer.err, true);
            addButton.textContent = 'Сохранить';
            addButton.disabled = false;
            return;
        }

        app.showPopup(
            `Заявление составлено, время ожидания до 3-х рабочих дней`,
            false
        );

        setTimeout(() => {
            pages.main();
            addButton.textContent = 'Добавить';
            addButton.disabled = false;
        }, 3000);
      },
      searchUserByBrCode: (answer) => {
        const addButton = document.getElementById('addButton');

        if (!answer.ok) {
            app.showPopup(answer.err, true);
            addButton.textContent = 'Сохранить';
            addButton.disabled = false;
            return;
        }

        app.showPopup(
            `Клиент найден!`,
            false
        );

        setTimeout(() => {
            pages.formDefect(answer.user);
            addButton.textContent = 'Добавить';
            addButton.disabled = false;
        }, 1500);
      }
    },
    
    searchUserByBrCode: (form) => {
        const addButton = document.getElementById('addButton');
        if (addButton) {
            addButton.disabled = true;
            addButton.textContent = "Отправка ...";
        }
        ajax.sendFormData(options.apiUrl, form, app.receiver.searchUserByBrCode);
    },
    editClientDefect: (form) => {
        const addButton = document.getElementById('addButton');
        if (addButton) {
            addButton.disabled = true;
            addButton.textContent = "Отправка ...";
        }
        ajax.sendFormData(options.apiUrl, form, app.receiver.editClientDefect);
    },
    
    updateCities: () => {
      app.ajax("cities", {}, (answer) => {
        storage.citiesMapIdIndex = answer.reduce((map, city, index) => {
            map[city.id] = index
            return map;
        }, {})
        storage.cities = answer.map(item => {
            item.search = `${item.name} (${item.region})`.toLowerCase()
            return item;
        })
      });
    },
    updateProducts: () => {
      app.ajax("products", {}, (answer) => {
        storage.products = answer;
      });
    },
    getList: () => {

    },
    addClient: (form) => {
      const emailInput = form.querySelector('input[name="email"]');
      if (emailInput.value.trim() === '') {
        const now = new Date();
        const timestamp = now.toISOString().replace(/[-:.]/g, '').slice(0, 15);
        const randomDigits = Math.floor(Math.random() * 1000000000);
        emailInput.value = `not${timestamp}${randomDigits}@notmail.ru`;
      }
      const addButton = document.getElementById('addButton');
      addButton.disabled = true;
      addButton.textContent = "Отправка...";
      ajax.sendFormData(options.apiUrl, form, app.receiver.add);
    },
  
    fillForm: () => {
      pages.formAdd();
      // Заполнение текстовых полей
    document.getElementById('family_name').value = 'Иванов';
    document.getElementById('name').value = 'Иван';
    document.getElementById('middle_name').value = 'Иванович';
    document.getElementById('phone').value = '+7 (999) 999-99-99';
    document.getElementById('email').value = 'ivanov2@example.com';
    document.getElementById('city_name').value = 'Москва';
    document.getElementById('city_id').value = '1'; // Пример ID города
    document.getElementById('birth_date').value = '01/01/1990';
    document.getElementById('old_product_name').value = 'Vaporesso';
    document.getElementById('product_code').value = 'BR02345';
    document.getElementById('passport_number').value = '1234567890';

    // Выбор пола
    document.getElementById('male').checked = true;


    },
    
    differenceProducts: text => {
        let existsProduct = storage.products
            .reduce((result, item) => {
                result[item.name] = true 
                return result;
                }, 
            {})
        return text
            .split("\n")
            .map(item => item.trim())
            .filter(item => item.length > 0 && !existsProduct[item])
    },
  
    restrictInput: (inputElement) => {
      const cursorPosition = inputElement.selectionStart;
  
      let newValue = inputElement.value
        .replace(/[^A-Za-zА-Яа-я\s-]/g, "")
        .replace(/(\s)\s+/g, "$1")
        .replace(/(\-)\-+/g, "$1");
  
      if (newValue.startsWith(" ")) {
        newValue = newValue.trimStart();
      }
  
      if (inputElement.value !== newValue) {
        inputElement.value = newValue;
        inputElement.setSelectionRange(cursorPosition, cursorPosition);
      }
    },
    registrationTG: () => {
        // открыть ссылку тг регистрации
        utils.openLink(`${storage.tgBot}?start=registration`)
    },
    handleBlur: (inputElement) => {
      inputElement.value = inputElement.value.trimEnd();
    },
    formatPhoneNumber: (inputElement) => {
      // Убираем все символы, кроме цифр
      let input = inputElement.value.replace(/\D/g, "");
  
      if (input.startsWith("9")) {
        input = "7" + input;
      }
      if (input.startsWith("7") || input.startsWith("8")) {
        if (input.startsWith("8")) {
          input = "7" + input.slice(1);
        }
        let formatted = "+7 ";
        if (input.length > 1) {
          formatted += "(" + input.slice(1, 4);
        }
        if (input.length > 4) {
          formatted += ") " + input.slice(4, 7);
        }
        if (input.length > 7) {
          formatted += "-" + input.slice(7, 9);
        }
        if (input.length > 9) {
          formatted += "-" + input.slice(9, 11);
        }
        inputElement.value = formatted;
      } else {
        inputElement.value = input;
      }
    },
    startPhoneNumber: (inputElement) => {
      if (!inputElement.value) {
        inputElement.value = "+7 ";
      }
    },


    formatDateInput: (input) => {
      let value = input.value.replace(/\D/g, "");
      if (value.length > 2) {
        value = value.slice(0, 2) + "/" + value.slice(2);
      }
      if (value.length > 5) {
        value = value.slice(0, 5) + "/" + value.slice(5);
      }
      if (value.length > 10) {
        value = value.slice(0, 10);
      }
      input.value = value;
    },
    
    validateDateInput: (input) => {
      const value = input.value;
      const isValid = app.validateDate(value);
      const ageError = document.getElementById("age-error");
    
      if (isValid) {
        input.style.borderColor = "";
    
        const [day, month, year] = value.split("/");
        const birthDate = new Date(year, month - 1, day);
        const today = new Date();
        const age = today.getFullYear() - birthDate.getFullYear();
        const isAdult =
          age > 18 ||
          (age === 18 &&
            today >=
              new Date(
                today.getFullYear(),
                birthDate.getMonth(),
                birthDate.getDate()
              ));
    
        if (!isAdult) {
          ageError.textContent = "Вам должно быть не менее 18 лет";
          ageError.style.display = "inline";
          input.style.borderColor = "red";
          setTimeout(() => {
            input.value = "";
            ageError.style.display = "none";
            input.style.borderColor = "";
          }, 3000);
        } else {
          ageError.style.display = "none";
        }
      } else {
        input.style.borderColor = "red";
        ageError.textContent = "Введите корректную дату";
        ageError.style.display = "inline";
        setTimeout(() => {
          input.value = "";
          ageError.style.display = "none";
          input.style.borderColor = "";
        }, 3000);
      }
    },
    
    validateDate: (dateString) => {
      if (!/^\d{2}\/\d{2}\/\d{4}$/.test(dateString)) {
        return false;
      }
      const [day, month, year] = dateString.split("/").map(Number);
      if (day < 1 || day > 31 || month < 1 || month > 12) {
        return false;
      }
      const currentYear = new Date().getFullYear();
      if (year < currentYear - 100 || year > currentYear) {
        return false;
      }
      const date = new Date(year, month - 1, day);
      return (
        date.getFullYear() === year &&
        date.getMonth() === month - 1 &&
        date.getDate() === day
      );
    },
    updateCityId: (inputElement, hiddenIdElement, datalistId, cities) => {
      const inputValue = inputElement.value.trim().toLowerCase();
      const datalist = document.getElementById(datalistId);
      let foundCity = null;
  
      if (cities && cities.length > 0) {
        foundCity = cities.find((city) => city.name.toLowerCase() === inputValue);
      }
  
      if (foundCity) {
        hiddenIdElement.value = foundCity.id;
      } else {
        const selectedOption = Array.from(datalist.options).find(
          (option) => option.value.toLowerCase() === inputValue
        );
        if (selectedOption) {
          hiddenIdElement.value = selectedOption.getAttribute("data-id");
        } else {
          hiddenIdElement.value = "";
        }
      }
    },
    validateCityInput: (inputElement, hiddenIdElement) => {
      const errorMessage = document.getElementById("city-error");
      if (!hiddenIdElement.value) {
        inputElement.value = "";
        errorMessage.style.display = "block";
      } else {
        errorMessage.style.display = "none";
      }
    },
    hideCitySuggestions: (event) => {
      if (!event.target.closest("#city_name")) {
        document.getElementById("citySelect").style.display = "none";
      }
    },
    formatDeviceInput: (input) => {
      const value = input.value.toUpperCase();
      const cleanedValue = value.replace(/[^A-Z0-9]/g, "");
      if (cleanedValue.length <= 2) {
        input.value = cleanedValue.replace(/[^A-Z]/g, "");
      } else {
        const letters = cleanedValue.slice(0, 2).replace(/[^A-Z]/g, "");
        const numbers = cleanedValue.slice(2).replace(/\D/g, "").slice(0, 5);
        input.value = letters + numbers;
      }
    },
    formatPassportNumInput: (input) => {
      input.value = input.value.replace(/\D/g, "");
    },
    updateFileNames:() => {
      const input = document.getElementById('passport_photos');
      const fileNames = Array.from(input.files).map(file => file.name);
      document.getElementById('file-names').textContent = fileNames.join(', ') || 'Нет загруженных файлов';
    },
    showPopup: (message, isError = false) => {
      const popupMessage = document.getElementById("popup-message");
      const closeBtn = document.querySelector("#closeSvg");
      const popupBox = document.querySelector("#popup");
      popupMessage.textContent = message;
      if (isError) {
        popupBox.classList.add("error");
        closeBtn.classList.add("error");
      } else {
        popupBox.classList.remove("error");
        closeBtn.classList.remove("error");
      }
      document.getElementById("popup").style.display = "flex";
    },
    closePopup: () => {
      document.getElementById("popup").style.display = "none";
    },
    
    // Аааа каша
    showSuggestions: (cities) => {
      if (cities.length == 0) return;
      
      const citySelect = document.getElementById("citySelect")
      if (!citySelect) return;
      
      //citySelect.innerHTML = "";
  
      if (Object.keys(cities).length === 0) {
        citySelect.style.display = "none";
        return;
      }
  
      citySelect.style.display = "block";
      citySelect.addEventListener("scroll", () => {
        document.getElementById("city_name").blur();
      });
  
      Object.values(cities).forEach((city) => {
        const cityOption = document.createElement("div");
        cityOption.textContent = `${city.name} (${city.region})`;
        cityOption.addEventListener("click", function () {
          document.getElementById(
            "city_name"
          ).value = `${city.name} (${city.region})`;
          document.getElementById("city_id").value = city.id;
          citySelect.style.display = "none";
          document.getElementById("city_name").focus();
          app.validateCityInput(
            document.getElementById("city_name"),
            document.getElementById("city_id")
          );
        });
        citySelect.appendChild(cityOption);
      });
    },
  
    handleCityInput______X_Old: (input, cities) => {
      const inputValue = input.value.toLowerCase();
      if (inputValue.length < 2) {
        document.getElementById("citySelect").style.display = "none";
        document.getElementById("city_id").value = "";
        return;
      }
  
      const filteredCities = Object.values(cities).filter((city) =>
        city.name.toLowerCase().startsWith(inputValue)
      );
  
      app.showSuggestions(filteredCities);
  
      const exactMatch = Object.values(cities).find(
        (city) => city.name.toLowerCase() === inputValue
      );
      if (exactMatch) {
        document.getElementById("city_id").value = exactMatch.id;
      } else {
        document.getElementById("city_id").value = "";
      }
    },
    
    // 
    handleCityInput: (input) => {
      const inputValue = input.value.trim().toLowerCase()
      const citySelect = document.getElementById("citySelect")
      if (!citySelect) return;
      
      if (inputValue.length < 2) {
        citySelect.style.display = "none"
        return;
      }
  
        const likeCities = storage.cities
        .filter(item => item.search.includes(inputValue))
        .sort((a, b) => {
            const aStarts = a.search.startsWith(inputValue)
            const bStarts = b.search.startsWith(inputValue)
            
            if (aStarts && !bStarts) return -1;
            if (!aStarts && bStarts) return 1;
            return 0;
        })

      if (likeCities == 0) {
        citySelect.style.display = "none"
        return;
      } else {
        citySelect.style.display = "block"
      }
      
      storage.selectedCityID = likeCities[0].id
      citySelect.innerHTML = view.listSelectCity(likeCities);
    },
    
    selectCity:(cityID) => {
        console.log('selectCity', cityID)
        storage.selectedCityID = cityID
        const city = storage.cities[storage.citiesMapIdIndex[cityID]];
        utils.writeInput('city_name', `${city.name} (${city.region})`)
        utils.writeInput('city_id', cityID)
        
        const citySelect = document.getElementById("citySelect")
        if (!citySelect) return;
        citySelect.style.display = "none"
    },
    
    
    
    
    cache: new Map(),
    abortController: null,
  
    debounce: (func, delay) => {
      let timeoutId;
      return function (...args) {
        clearTimeout(timeoutId);
        timeoutId = setTimeout(() => {
          func.apply(this, args);
        }, delay);
      };
    },
  
    fetchSuggestions: async () => {
      const apiKey = "9564001e-c30e-41f6-99cc-7f934923a3b6";
      if (state.address_shop.length < 3) {
        state.suggestions = [];
        app.renderSuggestions();
        return;
      }
  
      if (app.cache.has(state.address_shop)) {
        state.suggestions = app.cache.get(state.address_shop);
        app.renderSuggestions();
        return;
      }
  
      if (app.abortController) {
        app.abortController.abort();
      }
      app.abortController = new AbortController();
  
      state.suggestionLoading = true;
  
      try {
        const response = await fetch(
          `https://suggest-maps.yandex.ru/v1/suggest?apikey=${apiKey}&print_address=1&text=${encodeURIComponent(
            state.address_shop
          )}&lang=ru_RU`,
          { signal: app.abortController.signal }
        );
  
        if (!response.ok) {
          throw new Error("Ошибка сети: " + response.statusText);
        }
  
        const data = await response.json();
        state.suggestions = data.results || [];
        app.cache.set(state.address_shop, state.suggestions);
      } catch (error) {
        if (error.name !== "AbortError") {
          state.errors.address_shop = [
            ["Ошибка при получении подсказок: " + error.message],
          ];
          console.error("Ошибка при получении подсказок:", error);
        }
      } finally {
        state.suggestionLoading = false;
        app.renderSuggestions();
      }
    },
  
    renderSuggestions: () => {
      const suggestionsList = document.getElementById("suggestions");
      suggestionsList.innerHTML = "";
      suggestionsList.style.display = "none";
  
      if (state.suggestions.length > 0) {
        state.suggestions.forEach((suggestion) => {
          const li = document.createElement("div");
          const subtitleText = suggestion.subtitle?.text || "";
          const titleText = suggestion.title?.text || "";
          li.textContent = `${
            subtitleText === "" ? "" : subtitleText + ", "
          } ${titleText}`;
  
          li.addEventListener("click", () => {
            document.getElementById("address").value =
              suggestion.address.formatted_address;
            suggestionsList.innerHTML = "";
            suggestionsList.style.display = "none";
          });
          suggestionsList.appendChild(li);
          suggestionsList.style.display = "block";
        });
      }
    },
  
    handleAddressInput: function (event) {
      const input = event.target;
      const query = input.value.trim();
  
      state.address_shop = query;
  
      if (state.address_shop === "") {
        state.suggestions = [];
        app.renderSuggestions();
        return;
      }
  
      this.debounce(this.fetchSuggestions, 1000)();
    },
    showDeviceList: (deviceList) => {
      const deviceSelect = document.getElementById("deviceSelect");
      deviceSelect.innerHTML = "";
    
      if (Object.keys(deviceList).length === 0) {
        deviceSelect.style.display = "none";
        return;
      }
    
      deviceSelect.style.display = "block";
      deviceSelect.addEventListener("scroll", () => {
        document.getElementById("old_product_name").blur(); 
      });
    
      Object.values(deviceList).forEach((device) => {
        const deviceOption = document.createElement("div");
        deviceOption.textContent = device.name;
        deviceOption.addEventListener("click", function () {
          document.getElementById("old_product_name").value = device.name;
          document.getElementById("old_product_id").value = device.id;
          deviceSelect.style.display = "none";
          document.getElementById("old_product_name").focus();
        });
        deviceSelect.appendChild(deviceOption);
      });
    },
    
    handleDeviceInput: (input, deviceList) => {
      const inputValue = input.value.toLowerCase();
      if (inputValue.length < 2) {
        document.getElementById("deviceSelect").style.display = "none";
        document.getElementById("old_product_id").value = "";
        return;
      }
    
      const filteredList = Object.values(deviceList).filter((device) =>
        device.name.toLowerCase().includes(inputValue)
      );
    
      app.showDeviceList(filteredList);
    
      const exactMatch = Object.values(deviceList).find(
        (device) => device.name.toLowerCase() === inputValue
      );
      if (exactMatch) {
        document.getElementById("old_product_id").value = exactMatch.id;
      } else {
        document.getElementById("old_product_id").value = "";
      }
    },
    
    selectOldProduct: (select) => {
      const selectedValue = select.value;
        const oldProductNameInput = document.getElementById('old_product_name');

        if (selectedValue === "0") {
            oldProductNameInput.style.display = 'block';
            oldProductNameInput.value = '';
            oldProductNameInput.focus();
        } else if (selectedValue) {
            const selectedOption = select.options[select.selectedIndex];
            oldProductNameInput.value = selectedOption.text;
            oldProductNameInput.style.display = 'none';
        } else {
            oldProductNameInput.value = '';
            oldProductNameInput.style.display = 'none';
        }
    },
    toggleDropdown: () => {
      const options = document.getElementById('options');
      options.style.display = options.style.display === 'block' ? 'none' : 'block';
    },
  
    selectOption: (option) => {
      const selectedText = document.getElementById('selected-text');
      const oldProductNameInput = document.getElementById('old_product_name');
      const oldProductIdInput = document.getElementById('old_product_id'); // Скрытое поле
      const value = option.getAttribute('data-value');
  
      selectedText.textContent = option.textContent;
      oldProductIdInput.value = value;
  
      if (value === "0") {
        oldProductNameInput.style.display = 'block';
        oldProductNameInput.setAttribute('required', 'required');
        oldProductNameInput.value = '';
        oldProductNameInput.focus();
      } else if (value) {
        oldProductNameInput.style.display = 'none';
        oldProductNameInput.removeAttribute('required');
        oldProductNameInput.value = option.textContent;
      } else {
        oldProductNameInput.value = '';
        oldProductNameInput.style.display = 'none';
        oldProductNameInput.removeAttribute('required');
      }
  
      app.toggleDropdown();
    },
  
    initSelectDevice: () => {
      const customSelectElement = document.querySelector('.custom-select');
      const optionsElement = document.getElementById('options');
      if (!customSelectElement || !optionsElement) {
        return;
      }
      const options = document.querySelectorAll('.option');
      options.forEach(option => {
        option.addEventListener('click', () => app.selectOption(option));
      });
    
      document.addEventListener('click', (event) => {
        if (!customSelectElement.contains(event.target)) {
          optionsElement.style.display = 'none';
        }
      });
    },
    request: {
        all: () => {
          app.ajax('requestsList', {}, (r) => {
            if (r.ok) {
              pages.requests(r.data);
            }
          });
        },
        apply: (id) => {
          app.ajax('requestSetData', {
            request_id: id,
            is_accepted: 1,
            is_denided: 0
          }, r => {
            if (r.ok) {
              app.request.all();
            }
          });
        },
        denine: (id) => {
          app.ajax('requestSetData', {
            request_id: id,
            is_accepted: 0,
            is_denided: 1
          }, r => {
            if (r.ok) {
              app.request.all();
            }
          })
        },
    },

    
}


// error to log
{
window.addEventListener('error', function(event) {
    let data = {
        'userID': 0,
        'message': event.message,
        'file': event.filename,
        'line': event.lineno + ':' + event.colno,
        'userAgent': navigator.userAgent,
        'url':location.href,
        'timeAfterStartPage': 0
    }
    let formData = new FormData()
    for (let key in data)
        formData.append(key, data[key])
    fetch('/api/log/jsErrors.php', {
        method: 'POST',
        body: formData
    })
})
}

app.main()
