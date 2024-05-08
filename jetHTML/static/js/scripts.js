 
    function очиститьФорму(event) {
        event.preventDefault();
        event.stopPropagation();
        console.log(event.target.form);
        let секцииБлоков = document.querySelectorAll('[id$="-section"]');
        секцииБлоков.forEach(function(секция) {
            let префиксИдБлока = секция.id.split("-")[0];
            let блоки = секция.querySelectorAll(`[id^="${префиксИдБлока}-"]`);
            блоки.forEach(function(блок, index) {
                if (index > 0) {
                    блок.remove();
                }
            });
        });
        event.target.form.reset()
    }

    function Шаблон (имяШаблона, данныеШаблона) {          
        let шаблоны =   {
            "role_access":  `


    <fieldset class="" id="role_access-${данныеШаблона.УИД}">     
            <div class="строка">
                <div class="элемент отступ-внутр-10 id="role-${данныеШаблона.УИД}">
                    <!-- <label class="label">Роль</label> -->
                
                    <div class="выпадающий по-наведению">                      
                        
                            <button class="кнопка контур" aria-haspopup="true"  >
                                <span>Выбрать роль</span>
                                <span class="иконка is-small">
                                <i class="fas fa-angle-down" aria-hidden="true"></i>
                                </span>
                            </button>                
                        <div class="выпадающее-меню" id="dropdown-menu" role="menu">
                        <div class="выпадающий-контент">
                            <hr class="разделитель" />
                            <label class="checkbox выпадающий-элемент">
                                <input type="checkbox"  name="роль[${данныеШаблона.УИД}]" value="1"/>
                                Роль 1
                            </label>
                            <hr class="разделитель" />
                            <label class="checkbox выпадающий-элемент">
                                <input type="checkbox"  name="роль[${данныеШаблона.УИД}]" value="2"/>
                                Роль 2
                            </label>
                            <hr class="разделитель" />
                            <label class="checkbox выпадающий-элемент">
                                <input type="checkbox"  name="роль[${данныеШаблона.УИД}]" value="3"/>
                                Роль 2
                            </label>
                            <hr class="разделитель" />
                            <label class="checkbox выпадающий-элемент">
                                <input type="checkbox"  name="роль[${данныеШаблона.УИД}]" value="4"/>
                                Роль 2
                            </label>
                            <hr class="разделитель" />
                            <label class="checkbox выпадающий-элемент">
                                <input type="checkbox"  name="роль[${данныеШаблона.УИД}]" value="5"/>
                                роль 2
                            </label>
                            <hr class="разделитель" />
                        </div>
                        </div>
                    </div>
                </div> 

                <div class="элемент отступ-внутр-10" id="access-${данныеШаблона.УИД}">
                    <!-- <label class="label">Права доступа</label> -->
                    <div class="выпадающий по-наведению">
                        <div class="dropdown-trigger">
                        <button class="кнопка контур" aria-haspopup="true" aria-controls="dropdown-menu">
                            <span>Выбрать права доступа</span>
                            <span class="icon is-small">
                            <i class="fas fa-angle-down" aria-hidden="true"></i>
                            </span>
                        </button>
                        </div>
                        <div class="выпадающее-меню " id="dropdown-menu" role="menu">
                            <div class="выпадающий-контент">
                                    <hr class="разделитель" />
                                    <label class="checkbox выпадающий-элемент">
                                        <input type="checkbox" name="права_доступа[${данныеШаблона.УИД}]" value="1"/>
                                        права_доступа 1
                                    </label>
                                    <hr class="разделитель" />
                                    <label class="checkbox выпадающий-элемент">
                                        <input type="checkbox"  name="права_доступа[${данныеШаблона.УИД}]" value="2"/>
                                        права_доступа 2
                                    </label>
                                    <hr class="разделитель" />
                                    <label class="checkbox выпадающий-элемент">
                                        <input type="checkbox"  name="права_доступа[${данныеШаблона.УИД}]" value="3"/>
                                        права_доступа 3
                                    </label>
                                    <hr class="разделитель " />
                                    <label class="checkbox выпадающий-элемент">
                                        <input type="checkbox"  name="права_доступа[${данныеШаблона.УИД}]" value="4"/>
                                        права_доступа 4
                                    </label>
                                    <hr class="разделитель " />
                                    <label class="checkbox выпадающий-элемент">
                                        <input type="checkbox" name="права_доступа[${данныеШаблона.УИД}]" value="5"/>
                                        права_доступа 5
                                    </label>
                                    <hr class="разделитель" />
                            </div>
                        </div>
                    </div>
                </div>    


                
    <div class="строка">       
            <button onclick='добавитьБлок(event, {&#34;УИД&#34;:&#34;${данныеШаблона.УИД}&#34;,&#34;имяШаблона&#34;:&#34;${данныеШаблона.имяШаблона}&#34;} )' class="кнопка с-иконкой основной"> 
                <i class="fas fa-plus p-1"></i>
            </button>          
            <button onclick='удалитьБлок(event, {&#34;УИД&#34;:&#34;${данныеШаблона.УИД}&#34;,&#34;имяШаблона&#34;:&#34;${данныеШаблона.имяШаблона}&#34;} )' class="кнопка с-иконкой внимание">
            <i class="fas fa-minus p-1"></i>     
    </div>
              
            </div>          
    </fieldset>
`,
            "service_handler": `<fieldset id="service_handler-${данныеШаблона.УИД}" class="" oninit="пересчитатьОчередь">
    <div class="строка">       
            <div class="элемент отступ-внутр-10">
                <input style="max-width: 80px;" class="кнопка контур" type="number" name="очередь[${данныеШаблона.УИД}]" value="${+данныеШаблона.очередь+1}" id="order-service_handler-${данныеШаблона.УИД}" min="1"  max="1" onchange="переместитьБлок(event)">
            </div>
            <div  class="элемент строка отступ-внутр-10">
                <span class="icon подсказка отступ-внутр-10" data-tooltip="Нужно для того чтобы понять в какой сервис отправлять запрос для обработчки маршрута или действия, если обработчик не задан. Если задан обработчик, то скорей всего он зарегистрирован в СинКвике и соответсвует какому то Сервису">
                    <i class="fas fa-info-circle" ></i>
                </span>
                <div class="select">
                    <select class="кнопка контур" id="service-1" name="сервис[${данныеШаблона.УИД}]" class="" onchange="изменитьСписокОбработчиков(event, '${данныеШаблона.УИД}')">
                        <option disabled="disabled" value="" selected="selected">Выбрать сервис</option>
                        <option value="option1">Сервис 1Сервис 1Сервис 1Сервис 1</option>
                        <option value="option2">Сервис 2</option>
                        <option value="option3">Сервис 3</option>
                        <option value="option4">Сервис 4</option>
                    </select>
                </div>
            </div>
            <div class="элемент строка отступ-внутр-10">
                <span class="icon подсказка отступ-внутр-10" data-tooltip="Имя обработчика который существуе в Сервисе, если задан то поле Сервис можно не заполнять">
                    <i class="fas fa-info-circle" ></i>
                </span>
                <div class="select"> 
                     <select id="handler-1" class="кнопка контур" name="обработчик[${данныеШаблона.УИД}]">
                            <option disabled="disabled" value="" selected="selected">Выбрать обработчик</option>
                            <option value="option2">Обработчик 2</option>
                            <option value="option3">Обработчик 3</option>
                            <option value="option4">Обработчик 4</option>
                        </select>                 
                </div>            
            </div>
    <div class="строка">       
            <button onclick='добавитьБлок(event, {&#34;функции&#34;:[&#34;пересчитатьОчередь&#34;],&#34;УИД&#34;:&#34;${данныеШаблона.УИД}&#34;,&#34;имяШаблона&#34;:&#34;${данныеШаблона.имяШаблона}&#34;} )' class="кнопка с-иконкой основной"> 
                <i class="fas fa-plus p-1"></i>
            </button>          
            <button onclick='удалитьБлок(event, {&#34;УИД&#34;:&#34;${данныеШаблона.УИД}&#34;,&#34;имяШаблона&#34;:&#34;${данныеШаблона.имяШаблона}&#34;,&#34;функции&#34;:[&#34;пересчитатьОчередь&#34;]} )' class="кнопка с-иконкой внимание">
            <i class="fas fa-minus p-1"></i>     
    </div>
</div>
</fieldset>
`,
        }       
        return шаблоны[имяШаблона]
    }  


    function ПолучитьЗначениеИнпута(ИдИнпута){
        return document.getElementById(ИдИнпута).value
    }



    function НовыйИд(){
        return Date.now().toString(36)
    }
    
    let Функции = {
        "ПолучитьЗначениеИнпута": ПолучитьЗначениеИнпута,
        "пересчитатьОчередь": пересчитатьОчередь,
        "НовыйИд": НовыйИд,
        "переместитьБлок": переместитьБлок,
    }
    function переместитьБлок(event) {
        event.preventDefault();
        event.stopPropagation();

        let полеОчереди = event.target;
        let новаяОчередь = parseInt(полеОчереди.value);
        let блок = полеОчереди.closest('[id^="service_handler-"]');
        let префиксИдБлока = блок.id.split("-")[0];
        let родительскаяСекция = document.getElementById(`${префиксИдБлока}-section`);
        let блоки = родительскаяСекция.querySelectorAll(`[id^="${префиксИдБлока}-"]`);

        if (новаяОчередь < 1 || новаяОчередь > блоки.length) {
          
            return;
        }

        let целеваяПозиция = новаяОчередь - 1;
        let текущаяПозиция = Array.from(блоки).indexOf(блок);

        if (целеваяПозиция !== текущаяПозиция) {
            if (целеваяПозиция < текущаяПозиция) {
                родительскаяСекция.insertBefore(блок, блоки[целеваяПозиция]);
            } else {
                родительскаяСекция.insertBefore(блок, блоки[целеваяПозиция].nextSibling);
            }
        }

        пересчитатьОчередь(префиксИдБлока);
    }
    function пересчитатьОчередь(префиксИдБлока) {
        // let префиксИдБлока = текущийЭлемент.id.split("-")[0];

        let родительскаяСекция = document.getElementById(`${префиксИдБлока}-section`);
        let блоки = родительскаяСекция.querySelectorAll(`[id^="${префиксИдБлока}-"]`);
        let максимальнаяОчередь = блоки.length;
        блоки.forEach(function(блок, index) {
            let полеОчереди = блок.querySelector('input[name^="очередь"]');
            полеОчереди.value = index + 1;
            полеОчереди.max=максимальнаяОчередь
        });
    }

    // function добавитьБлок(event, имяШаблона, УИД) {  
    function добавитьБлок(event, данные) {  
        event.preventDefault();
        event.stopPropagation();
 
        let блокКоторыйВызвалСобытие = event.target.closest(`[id^="${данные.имяШаблона}"]`);
        let префиксИдБлока = блокКоторыйВызвалСобытие.id.split("-")[0];

        let родительскийБлок = блокКоторыйВызвалСобытие.closest(`[id^="${префиксИдБлока}-section"]`);
        
        let ИдНовогоБлока = НовыйИд();

        данные["УИД"] = ИдНовогоБлока;
        let новыйБлок = Шаблон(данные.имяШаблона,  данные);
        блокКоторыйВызвалСобытие.insertAdjacentHTML('afterend', новыйБлок);

        // приверим если в данных кнопки был массиф функций, то вызовим их, передав в качестве аргумента блок вызвавший текущее событие
        if (данные.hasOwnProperty("функции")) {
            функцииПостОбраотки  = данные["функции"];
            функцииПостОбраотки.forEach(function(функция) {
                Функции[функция](данные.имяШаблона);
            })
        }      
    }

    // function удалитьБлок(event, имяШаблона, УИД) {
    function удалитьБлок(event, данные) {
        event.preventDefault();
        event.stopPropagation();
        let родительскаяСекция = event.target.closest(`[id="${данные.имяШаблона}-section"]`);
        // если был удалён последний элемент в сексии блоков, то добавим новый
      
       
        let блокКоторыйВызвалСобытие = event.target.closest(`[id^="${данные.имяШаблона}"]`);
            блокКоторыйВызвалСобытие.remove();
        let блоки = document.querySelectorAll(`[id^="${данные.имяШаблона}"]`);

    

        if (блоки.length == 0) {
            данные["УИД"] = НовыйИд();
            let новыйБлок = Шаблон(данные.имяШаблона, данные);
          
            // найдём родительскую секцию и вставим внеё новый блок
            родительскаяСекция.insertAdjacentHTML('afterbegin', новыйБлок);
            // document.getElementById(имяШаблона+"-section").insertAdjacentHTML('afterbegin', новыйБлок);
        }
        if (данные.hasOwnProperty("функции")) {
            функцииПостОбраотки  = данные["функции"];
            функцииПостОбраотки.forEach(function(функция) {
                Функции[функция](данные.имяШаблона);
            })
        }      
       
        
    }

    // Нужно добавить в рендер функцию, точнее в парсинг функцию обработчик который будет собирать все javascript из шаблонов и помещать их в один файл,  
    /*
    может можно сделать так чтобы каждый js блок кода был объявлен как define "обработчикФормыНовогоМаршрута.js"
    */

    var количествоОбработчиков = 0 // просто УИД добавленного обработчика, не уменьшается при удалении, потому 


