
/*
Обработка нажатия кнопок ajax ? в соновном меню, личный кабинет, не требующих открытия новой вкладки
*/ 

function ajax(event, изменитьАдреснуюСтроку = false) {
  event.preventDefault();
  event.stopPropagation();  
  let target = event.currentTarget   
  console.log("target", target);
  let действие = target.getAttribute('href');
  // const headers = response.headers;
  
  // let данныеФормы = new FormData(target); // создаем объект FormData и автоматически парсим форму
  // данныеФормы.append("действие", "добавитьОбработчик")
  if (изменитьАдреснуюСтроку) {
    console.log("действие", действие);
    history.pushState({}, '', действие);

  }
  fetch(действие, {
        method: 'AJAX',   
        headers: {
          'Method': 'AJAX'
        }
      }).then(response => response.json())
        .then(данные => {
          console.log(данные);
          ОбработатьОтветСервера(данные)
          // for ( имяБлока in data) {

          //   // Обработчик вставка, нужно вынести в отдельную функцию
          //   console.log(data[имяБлока]);
          //   let данныеБлока = data[имяБлока]
          //   if (!!данныеБлока.HTML && данныеБлока.HTML !== "") {            
          //     let template = document.createElement("template");
          //     template.innerHTML = данныеБлока.HTML               
          //     let первыйЭлемент = template.content.firstElementChild
          //     if (!!первыйЭлемент.dataset){
          //       let методВставки  = первыйЭлемент.dataset.updatemethod                  
          //       let id = первыйЭлемент.id  

          //       console.log("первыйЭлемент", первыйЭлемент, "методВставки", методВставки, "id", id);
                
          //       document.getElementById(id)[методВставки](template.content)
          //     }
          //   } 
          // }
      }).catch(error => {
          console.log(error);
      });
}

function ajaxPost(event){
  event.preventDefault();
  event.stopPropagation();  
  let форма = event.target;  
  if (форма.tagName!=="FORM" )  {
    console.log("не форма ", форма);

    форма = форма.form

    if (форма.tagName!=="FORM" )  {
      console.log("всё ещё не форма ", форма);
    }
  }
  console.log("event", event);
    let действие = форма.getAttribute('action');
    console.log("действие формы", действие);
  let данныеФорм;
  // if (форма.hasAttribute('beforeSubmit')) {
  //   обработчикПередОбтправкой = форма.getAttribute('beforeSubmit');
  //   данныеФормы = new FormData()
  //   данныеФормы.append("действие", действие)
  //   console.log(обработчикПередОбтправкой, Функции[обработчикПередОбтправкой]);
  //   let структурированныеДанные = Функции[обработчикПередОбтправкой](event, форма)
  //   console.log(структурированныеДанные);
  //   данныеФормы.append("данные", JSON.stringify(структурированныеДанные))

  // } else {

// }
    данныеФорм = new FormData(форма); // создаем объект FormData и автоматически парсим форму
    данныеФорм.append("действие", действие)
    console.groupCollapsed(`FormData ${new Date().toLocaleTimeString()}`);
    for (var pair of данныеФорм.entries()) {
        console.log(`${pair[0]} : ${pair[1]}`);
    }
    console.groupEnd();


  fetch(`/${действие}`, {
        method: 'AJAXPost',
        body: данныеФорм,
        headers: {
          'Method': 'AJAXPost'
        }
      }).then(response => response.json())
        .then(данные => {
          console.log(данные);
          ОбработатьОтветСервера(данные)
          // for ( имяБлока in data) {

          //   // Обработчик вставка ньд нужно вынести в отдельную функцию
          //   console.log(data[имяБлока]);
          //   let данныеБлока = data[имяБлока]
          //   if (!!данныеБлока.HTML && данныеБлока.HTML !== "") {       

          //     let HTML = document.createElement("template");
          //         HTML.innerHTML = данныеБлока.HTML    

          //     let первыйЭлемент = кодВставки.content.firstElementChild // находим первый элемент вставки
            
          //     if (!!первыйЭлемент.dataset){
              
                            
          //       let селекторКонтейнера = первыйЭлемент.getAttribute("data-parent_selector"); // селектор элемента в который нужно вставить полученные данные 

          //       if (!селекторКонтейнера) {                  
          //         селекторКонтейнера = первыйЭлемент.id         // если нет data-parent_selector то берём id первого элемента из вставлемого html                      
          //       }

          //       const контенйреДляВставки = document.querySelector(селекторКонтейнера); // находим контейнер куда будет вставлен полученные данные
          //       if (!!селекторКонтейнера) {
          //         let методВставки = первыйЭлемент.getAttribute("updatemethod"); // получим метод вставки из перовго элемента блока вставки, если такого нет то проверим у родителя в который будем вставлять
          //         if (!методВставки) {
          //           методВставки = контенйреДляВставки.getAttribute("updatemethod");
          //         }
          //         if (!!методВставки){
          //           контенйреДляВставки[методВставки](HTML.content)
          //         } else {
          //           console.warn("Не найден метод обновления контента для вставки ?? Перезаписывать ?");
          //           контенйреДляВставки.innerHTML= HTML.content;
          //         }    
          //        } else {
          //         console.warn("Не найден контейнер для вставки");
          //         ВсплывающееСообщение("Не найден контейнер для вставки", тип="внимание", время=0)
          //       }
          //     }
          //   }
          // }
      }).catch(error => {
          console.log(error);
      });
}


function ОбработатьОтветСервера(данные){
  for ( имяБлока in данные) {

    console.log(имяБлока, данные[имяБлока]);

    let данныеБлока = данные[имяБлока]
    if (!!данныеБлока.HTML && данныеБлока.HTML !== "") {       

      let HTML = document.createElement("template");
          HTML.innerHTML = данныеБлока.HTML    

      let первыйЭлемент = HTML.content.firstElementChild // находим первый элемент вставки
    
      if (!!первыйЭлемент.dataset){
      
                    
        let селекторКонтейнера = первыйЭлемент.getAttribute("data-parent_selector"); // селектор элемента в который нужно вставить полученные данные 

        let контейнерДляВставки = document.querySelector(селекторКонтейнера); // находим контейнер куда будет вставлен полученные данные
        if (!селекторКонтейнера) {                  
          селекторКонтейнера = первыйЭлемент.id         // если нет data-parent_selector то берём id первого элемента из вставлемого html    
          контейнерДляВставки =  document.getElementById(селекторКонтейнера)                
        }
        

      
        console.log(контейнерДляВставки, селекторКонтейнера)
        if (!!селекторКонтейнера) {
          // let методВставки = первыйЭлемент.getAttribute("updatemethod"); // получим метод вставки из перовго элемента блока вставки, если такого нет то проверим у родителя в который будем вставлять
          let методВставки = первыйЭлемент.dataset.updateMethod; // получим метод вставки из перовго элемента блока вставки, если такого нет то проверим у родителя в который будем вставлять
                 
         
          if (!методВставки) {
            // методВставки = контенйреДляВставки.getAttribute("updatemethod");
            методВставки = контейнерДляВставки.dataset.updateMethod;
          }

          if (!!методВставки){
            контейнерДляВставки[методВставки](HTML.content)
          } else {
            console.warn("Не найден метод обновления контента для вставки ?? Перезаписывать ?");
            контейнерДляВставки.replaceWith(HTML.content);
          }    
          ВсплывающееСообщение("Всё ок", тип="инфо", время=0)
          // ВсплывающееСообщение("Всё ок", тип="важно", время=0)
          // ВсплывающееСообщение("Всё ок", тип="внимание", время=0)
         } else {
          console.warn("Не найден контейнер для вставки");
          ВсплывающееСообщение("Не найден контейнер для вставки", тип="внимание", время=0)
        }
      }
    }
  }
}

const постоянныеШаблоны = {
  иконки:{
    "инфо":'<i class="fa-regular fa-lightbulb"></i>',
    "внимание":'<i class="fa-solid fa-bomb"></i>',
    "важно":'<i class="fa-solid fa-circle-exclamation"></i>',
    
  }
};

function ВсплывающееСообщение(текст, тип="инфо", время=0) {
  const контейнерСообщений = document.querySelector("#сообщения");
  let сообщение = document.createElement("div");
      сообщение.classList.add("сообщение","контур", "новое", тип);
      // сообщение.classList.add(тип);
      const иконка =постоянныеШаблоны.иконки[тип]

      сообщение.innerHTML = `<div class="иконка">${иконка}</div><div class="текст">${текст}</div>`;
      контейнерСообщений.append(сообщение);

      setTimeout(function() {
        сообщение.classList.remove("новое");
       
      }, 300);

  if (время > 0){
    setTimeout(function() { 
      setTimeout(function() {
        сообщение.classList.add("спрятать");
        setTimeout(function() {
          сообщение.remove();
        }, 300); // Длительность анимации
      }, время);
    }, время);
  }
  сообщение.addEventListener("click", function() {
    сообщение.classList.add("спрятать");
   
     setTimeout(function() {     
      // сообщение.remove();  
      // контейнерСообщений.classList.add("падение-вниз");
      setTimeout(function() {
        
        // контейнерСообщений.classList.remove("падение-вниз");
      }, 1300);  
    
    }, 500); // Длительность анимации

    });
}

/* ОткрытьСтраницу - Открывает новую вкладку и загружает страницу. для товаров к примеру. */
function ОткрытьСтраницу(event) { 
  event.preventDefault();
  event.stopPropagation();  
  let текущийЭлемент = event.currentTarget   
  let href = текущийЭлемент.dataset.href;    
  window.open(href, "_blank");  
}



function connectToWebSocketServer() {

  // Прямое соединение с RenderServerom для разработки, при изменении css, js, html обновляем страницу
    const socket = new WebSocket('ws://localhost:444');
  
 
    socket.onopen = function() {
      console.log('WebSocket соединение установленно');
    };
     // Event listener for when a message is received from the server
    socket.onmessage = function(event) {
      console.log('Получено сообщение с сервера:', event);
      if ( event.data == "reload"){
        socket.close();
        setTimeout(function() { location.reload() }, 1000);
       
      }
    };
   
  
    socket.onclose = function(event) {
      console.log('WebSocket cсоединение закрыто:', event.code, event.reason);
    };
  
    
    socket.onerror = function(error) {
      console.error('WebSocket ошибка:', error);
    };
  }
  connectToWebSocketServer() 

  document.addEventListener("submit", function(event) {
    console.log(event);
    // Все отправки форм отправляем ajax методом, перехватываем всё, чтобы не писать н акаждой форме onsubmit 
    ajaxPost(event);
  });
