 


    function Шаблон (имяШаблона, данныеШаблона) {          
        let шаблоны =   {
            "role_access":  `




    <fieldset class="" id="role_access-${данныеШаблона.УИД}">     
            <div class="строка">
                <div class="элемент отступ-внутр-10 id="role-${данныеШаблона.УИД}">
                    <!-- <label class="label">Роль</label> -->
                    <div class="выпадающий по-наведению">
                        <div class="dropdown-trigger">
                        <button class="кнопка контур" aria-haspopup="true" aria-controls="dropdown-menu">
                            <span>Выбрать роль</span>
                            <span class="иконка is-small">
                            <i class="fas fa-angle-down" aria-hidden="true"></i>
                            </span>
                        </button>
                        </div>
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
      
            <button onclick='добавитьБлок(event, {&#34;имяШаблона&#34;:&#34;${данныеШаблона.имяШаблона}&#34;,&#34;УИД&#34;:&#34;${данныеШаблона.УИД}&#34;} )' class="кнопка с-иконкой основной"> 
                <i class="fas fa-plus p-1"></i>
            </button>          
            <button onclick='удалитьБлок(event, {&#34;УИД&#34;:&#34;${данныеШаблона.УИД}&#34;,&#34;имяШаблона&#34;:&#34;${данныеШаблона.имяШаблона}&#34;} )' class="кнопка с-иконкой внимание">
            <i class="fas fa-minus p-1"></i>
        </button> <button onclick='удалитьБлок(event, {&#34;УИД&#34;:&#34;${данныеШаблона.УИД}&#34;,&#34;имяШаблона&#34;:&#34;${данныеШаблона.имяШаблона}&#34;} )' class="кнопка с-иконкой внимание">
            <i class="fas fa-minus p-1"></i>
        </button> <button onclick='удалитьБлок(event, {&#34;УИД&#34;:&#34;${данныеШаблона.УИД}&#34;,&#34;имяШаблона&#34;:&#34;${данныеШаблона.имяШаблона}&#34;} )' class="кнопка  с-иконкой внимание">
            <i class="fas fa-minus p-1"></i>
        </button> <button onclick='удалитьБлок(event, {&#34;УИД&#34;:&#34;${данныеШаблона.УИД}&#34;,&#34;имяШаблона&#34;:&#34;${данныеШаблона.имяШаблона}&#34;} )' class="кнопка с-иконкой внимание">
            <i class="fas fa-minus p-1"></i>
        </button>
       
    </div>

              
            </div>          
    </fieldset>
`,
            "service_handler": `<fieldset 1 id="service_handler-${данныеШаблона.УИД}" class="field is-horizontal">
    <div class="field-body">
       
            <div class="field control" style="max-width: 80px;">
                <input class="input" type="number" name="очередь[${данныеШаблона.УИД}]" value="${+данныеШаблона.очередь+1}" id="order-service_handler-${данныеШаблона.УИД}" min="1">
            </div>
            <div
                class="field control">
                <!-- <input class="input" type="text" placeholder="Имя сервиса" name ="сервис" value=""> -->
                <div class="select">
                    <select id="service-1" name="сервис[${данныеШаблона.УИД}]" class="" onchange="изменитьСписокОбработчиков(event, '${данныеШаблона.УИД}')">
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
                        <select id="handler-1" name="обработчик[${данныеШаблона.УИД}]">
                            <option disabled="disabled" value="" selected="selected">Выбрать обработчик</option>
                            <option value="option2">Обработчик 2</option>
                            <option value="option3">Обработчик 3</option>
                            <option value="option4">Обработчик 4</option>
                        </select>
                    </div>
                </div>
            </div>
           
            

    <div class="строка"> 
      
            <button onclick='добавитьБлок(event, {&#34;имяШаблона&#34;:&#34;${данныеШаблона.имяШаблона}&#34;,&#34;очередь&#34;:&#34;${+данныеШаблона.очередь+1}&#34;,&#34;УИД&#34;:&#34;${данныеШаблона.УИД}&#34;} )' class="кнопка с-иконкой основной"> 
                <i class="fas fa-plus p-1"></i>
            </button>          
            <button onclick='удалитьБлок(event, {&#34;очередь&#34;:&#34;${+данныеШаблона.очередь+1}&#34;,&#34;УИД&#34;:&#34;${данныеШаблона.УИД}&#34;,&#34;имяШаблона&#34;:&#34;${данныеШаблона.имяШаблона}&#34;} )' class="кнопка с-иконкой внимание">
            <i class="fas fa-minus p-1"></i>
        </button> <button onclick='удалитьБлок(event, {&#34;очередь&#34;:&#34;${+данныеШаблона.очередь+1}&#34;,&#34;УИД&#34;:&#34;${данныеШаблона.УИД}&#34;,&#34;имяШаблона&#34;:&#34;${данныеШаблона.имяШаблона}&#34;} )' class="кнопка с-иконкой внимание">
            <i class="fas fa-minus p-1"></i>
        </button> <button onclick='удалитьБлок(event, {&#34;УИД&#34;:&#34;${данныеШаблона.УИД}&#34;,&#34;имяШаблона&#34;:&#34;${данныеШаблона.имяШаблона}&#34;,&#34;очередь&#34;:&#34;${+данныеШаблона.очередь+1}&#34;} )' class="кнопка  с-иконкой внимание">
            <i class="fas fa-minus p-1"></i>
        </button> <button onclick='удалитьБлок(event, {&#34;очередь&#34;:&#34;${+данныеШаблона.очередь+1}&#34;,&#34;УИД&#34;:&#34;${данныеШаблона.УИД}&#34;,&#34;имяШаблона&#34;:&#34;${данныеШаблона.имяШаблона}&#34;} )' class="кнопка с-иконкой внимание">
            <i class="fas fa-minus p-1"></i>
        </button>
       
    </div>


    </div>


</fieldset>

`,
        }       
        // console.log(шаблоны[имяШаблона]);
        return шаблоны[имяШаблона]
    }  



    function НовыйИд(){
        return Date.now().toString(36)
    }

    // function добавитьБлок(event, имяШаблона, УИД) {  
    function добавитьБлок(event, данные) {  
        event.preventDefault();
        event.stopPropagation();
        console.log(данные);
        console.log(`[id^="${данные.имяШаблона}"]`);

        let блокКоторыйВызвалСобытие = event.target.closest(`[id^="${данные.имяШаблона}"]`);


        let ИдНовогоБлока = НовыйИд();
        // let новыйИд = `${данные.имяШаблона}-${ИдНовогоБлока}`;

        console.log("данные", данные, ИдНовогоБлока, блокКоторыйВызвалСобытие);

        данные["УИД"] = ИдНовогоБлока;
        let новыйБлок = Шаблон(данные.имяШаблона,  данные);
       console.log(новыйБлок);
        // console.log("блокКоторыйВызвалСобытие",`[id^="${данные.имяШаблона}"]`,  блокКоторыйВызвалСобытие);

        блокКоторыйВызвалСобытие.insertAdjacentHTML('afterend', новыйБлок);
    }

    // function удалитьБлок(event, имяШаблона, УИД) {
    function удалитьБлок(event, данные) {
        event.preventDefault();
        event.stopPropagation();
        let родительскаяСекция = event.target.closest(`[id="section-${данные.имяШаблона}"]`);
        // если был удалён последний элемент в сексии блоков, то добавим новый
      
       
        let блокКоторыйВызвалСобытие = event.target.closest(`[id^="${данные.имяШаблона}"]`);
            блокКоторыйВызвалСобытие.remove();
        let блоки = document.querySelectorAll(`[id^="${данные.имяШаблона}"]`);

        console.log(блоки);

        if (блоки.length == 0) {
            данные["УИД"] = НовыйИд();
            let новыйБлок = Шаблон(данные.имяШаблона, данные);
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

    var количествоОбработчиков = 0 // просто УИД добавленного обработчика, не уменьшается при удалении, потому 


