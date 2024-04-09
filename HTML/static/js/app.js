
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
    history.pushState({}, '', действие);

  }
  fetch(действие, {
        method: 'AJAX',   
        headers: {
          'Method': 'AJAX'
        }
      }).then(response => response.json())
        .then(data => {
          console.log(data);
          for ( имяБлока in data) {

            // Обработчик вставка ньд нужно вынести в отдельную функцию
            console.log(data[имяБлока]);
            let данныеБлока = data[имяБлока]
            if (!!данныеБлока.HTML && данныеБлока.HTML !== "") {            
              let template = document.createElement("template");
              template.innerHTML = данныеБлока.HTML               
              let первыйЭлемент = template.content.firstElementChild
              if (!!первыйЭлемент.dataset){
                let методВставки  = первыйЭлемент.dataset.updatemethod                  
                let id = первыйЭлемент.id  
                document.getElementById(id)[методВставки](template.content)
              }
            }
          }
      }).catch(error => {
          console.log(error);
      });
}

function ajaxPost(event){
  event.preventDefault();
  event.stopPropagation();  
  let target = event.target;   
  console.log("event", event);
  let действие = target.getAttribute('action');
  let formData = new FormData(target); // создаем объект FormData и автоматически парсим форму
  formData.append("действие", действие)
  fetch(действие, {
        method: 'AJAXPost',
        body: formData,
        headers: {
          'Method': 'AJAXPost'
        }
      }).then(response => response.json())
        .then(data => {
          console.log(data);
          for ( имяБлока in data) {

            // Обработчик вставка ньд нужно вынести в отдельную функцию
            console.log(data[имяБлока]);
            let данныеБлока = data[имяБлока]
            if (!!данныеБлока.HTML && данныеБлока.HTML !== "") {            
              let template = document.createElement("template");
              template.innerHTML = данныеБлока.HTML               
              let первыйЭлемент = template.content.firstElementChild
              if (!!первыйЭлемент.dataset){
                let методВставки  = первыйЭлемент.dataset.updatemethod                  
                let id = первыйЭлемент.id  
                document.getElementById(id)[методВставки](template.content)
              }
            }
          }
      }).catch(error => {
          console.log(error);
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
    const socket = new WebSocket('wss://localhost:444');
  
 
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
    // Все отправки форм отправляем ajax методом, перехватываем всё, чтобы не писать н акаждой форме onsubmit 
    ajaxPost(event);
  });