 

   
    function Шаблон (имяШаблона, номер) {   
       
        console.log("Шаблон", имяШаблона, номер);
        let шаблоны =   {
            "role_access": `<fieldset class="field is-horizontal" id="role_access-${номер}">     
            <div class="field-body">
                <div class="control field" id="role-1">
                    
                    <div class="dropdown is-hoverable">
                        <div class="dropdown-trigger">
                        <button class="button" aria-haspopup="true" aria-controls="dropdown-menu">
                            <span>Выбрать роль</span>
                            <span class="icon is-small">
                            <i class="fas fa-angle-down" aria-hidden="true"></i>
                            </span>
                        </button>
                        </div>
                        <div class="dropdown-menu" id="dropdown-menu" role="menu">
                        <div class="dropdown-content">
                        <div class="dropdown-item">
                            <label class="checkbox ">
                                <input type="checkbox"  name="роль[${номер}]" value="1"/>
                                Роль 1
                            </label>
                            <hr class="dropdown-divider" />
                            <label class="checkbox">
                                <input type="checkbox"  name="роль[${номер}]" value="2"/>
                                Роль 2
                            </label>
                            <hr class="dropdown-divider" />
                            <label class="checkbox">
                                <input type="checkbox"  name="роль[${номер}]" value="3"/>
                                Роль 2
                            </label>
                            <hr class="dropdown-divider" />
                            <label class="checkbox">
                                <input type="checkbox"  name="роль[${номер}]" value="4"/>
                                Роль 2
                            </label>
                            <hr class="dropdown-divider" />
                            <label class="checkbox">
                                <input type="checkbox"  name="роль[${номер}]" value="5"/>
                                роль 2
                            </label>
                        </div>
                        </div>
                        </div>
                    </div>
                </div> 

                <div class="field control" id="access-${номер}">
                    
                    <div class="dropdown is-hoverable">
                        <div class="dropdown-trigger">
                        <button class="button" aria-haspopup="true" aria-controls="dropdown-menu">
                            <span>Выбрать права доступа</span>
                            <span class="icon is-small">
                            <i class="fas fa-angle-down" aria-hidden="true"></i>
                            </span>
                        </button>
                        </div>
                        <div class="dropdown-menu" id="dropdown-menu" role="menu">
                        <div class="dropdown-content">
                        <div class="dropdown-item">
                            <label class="checkbox ">
                                <input type="checkbox" name="права_доступа[${номер}]" value="1"/>
                                права_доступа 1
                            </label>
                            <hr class="dropdown-divider" />
                            <label class="checkbox">
                                <input type="checkbox"  name="права_доступа[${номер}]" value="2"/>
                                права_доступа 2
                            </label>
                            <hr class="dropdown-divider" />
                            <label class="checkbox">
                                <input type="checkbox"  name="права_доступа[${номер}]" value="3"/>
                                права_доступа 3
                            </label>
                            <hr class="dropdown-divider" />
                            <label class="checkbox">
                                <input type="checkbox"  name="права_доступа[${номер}]" value="4"/>
                                права_доступа 4
                            </label>
                            <hr class="dropdown-divider" />
                            <label class="checkbox">
                                <input type="checkbox" name="права_доступа[${номер}]" value="5"/>
                                права_доступа 5
                            </label>
                        </div>
                        </div>
                        </div>
                    </div>
                </div>    
                <div class="field control has-addons">
        <div class="control ">
            <button onclick="добавитьБлок(event, '${имяШаблона}', '${номер}')" class="button  is-primary control"> 
                <i class="fas fa-plus p-1"></i>
            </button>
                    </div>
        <div class="control ">
            <button onclick="удалитьБлок(event, '${имяШаблона}', &#34;${номер}&#34;)" class="button  is-danger control">
            <i class="fas fa-minus p-1"></i>
        </button>
        </div>
    </div>
              
            </div>          
    </fieldset>
`,
            "service_handler": `<fieldset id="service_handler-${номер}" class="field is-horizontal">
    <div class="field-body">
            <div class="field control" style="max-width: 80px;">
                <input class="input" type="number" name="order[${номер}]" value="" id="order-service_handler-${номер}" min="1">
            </div>
            <div
                class="field control">
                
                <div class="select">
                    <select id="service-1" name="сервис[${номер}]" class="" onchange="изменитьСписокОбработчиков(event, &#34;${номер}&#34;)">
                        <option disabled="disabled" value="" selected="selected">Выбрать сервис</option>
                        <option value="option1">Сервис 1Сервис 1Сервис 1Сервис 1</option>
                        <option value="option2">Сервис 2</option>
                        <option value="option3">Сервис 3</option>
                        <option value="option4">Сервис 4</option>
                    </select>
                </div>
            </div>
            <div class="field control">
                <div class="select">
                    <div class="control">
                        <select id="handler-1" name="обработчик[${номер}]">
                            <option disabled="disabled" value="" selected="selected">Выбрать обработчик</option>
                            <option value="option2">Обработчик 2</option>
                            <option value="option3">Обработчик 3</option>
                            <option value="option4">Обработчик 4</option>
                        </select>
                    </div>
                </div>
            </div>
            <div class="field control has-addons">
        <div class="control ">
            <button onclick="добавитьБлок(event, '${имяШаблона}', '${номер}')" class="button  is-primary control"> 
                <i class="fas fa-plus p-1"></i>
            </button>
                    </div>
        <div class="control ">
            <button onclick="удалитьБлок(event, '${имяШаблона}', &#34;${номер}&#34;)" class="button  is-danger control">
            <i class="fas fa-minus p-1"></i>
        </button>
        </div>
    </div>

    </div>


</fieldset>

`,
        }       
        console.log(шаблоны[имяШаблона]);
        return шаблоны[имяШаблона]
    }  

    function добавитьБлок(event, имяШаблона, номер) {  
        event.preventDefault();
        event.stopPropagation();

        let блоки = document.querySelectorAll(`[id^="${имяШаблона}"]`);
        console.log(блоки);
// теперь elements содержит все элементы с id, начинающимся с "yourString"

        let количествоБлоков = блоки.length;
        let номерНовогоБлока = +количествоБлоков + 1;
        let новыйИд = `${имяШаблона}-${номерНовогоБлока}`;
        console.log(имяШаблона, номерНовогоБлока);
        let новыйБлок = Шаблон(имяШаблона, номерНовогоБлока);
        let блокКоторыйВызвалСобытие = event.target.closest(`[id^="${имяШаблона}"]`);
        console.log("блокКоторыйВызвалСобытие", блокКоторыйВызвалСобытие);

        блокКоторыйВызвалСобытие.insertAdjacentHTML('afterend', новыйБлок);
    }

    function удалитьБлок(event, имяШаблона, номер) {
        event.preventDefault();
        event.stopPropagation();
        let родительскаяСекция = event.target.closest(`[id="section-${имяШаблона}"]`);
        // если был удалён последний элемент в сексии блоков, то добавим новый
      
       
        let блокКоторыйВызвалСобытие = event.target.closest(`[id^="${имяШаблона}"]`);
            блокКоторыйВызвалСобытие.remove();
        let блоки = document.querySelectorAll(`[id^="${имяШаблона}"]`);

        console.log(блоки);

        if (блоки.length == 0) {
            let новыйБлок = Шаблон(имяШаблона, 1);
            console.log(родительскаяСекция);
            // найдём родительскую секцию и вставим внеё новый блок
            родительскаяСекция.insertAdjacentHTML('afterbegin', новыйБлок);
            // document.getElementById(имяШаблона+"-section").insertAdjacentHTML('afterbegin', новыйБлок);
        }
        
    }

    // Нужно добавить в рендер функцию, точнее в парсинг функцию обработчик который будет собирать все javascript из шаблонов и помещать их в один файл,  
    /*
    может можно сделать так чтобы каждый js блок кода был объявлен как define "обработчикФормыНовогоМаршрута.js"
    */

    var количествоОбработчиков = 0 // просто номер добавленного обработчика, не уменьшается при удалении, потому 


