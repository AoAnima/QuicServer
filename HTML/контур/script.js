const кнопкаОчисткиСтен = document.getElementById('оичтистьХолст');
const кнопкаРисованияДверей = document.getElementById('рисоватьДверь');
const холст = document.getElementById('холст');
const контекст = холст.getContext('2d');
const полеШагаСетки = document.getElementById('шагСетки');
const полеТолщиныСтены = document.getElementById('толщинаСтен');
const кнопкаРисованияСтен = document.getElementById('рисоватьСтену');
const полеДлиныСтены = document.getElementById('длинаСтены');
const полеШириныДвери = document.getElementById('ширинаДвери');


let режимРисованияСтен = false;
let режимРисованияДверей = false;
let началоX, началоY;
let длинаСтены = 0;

let комнаты = [];
let стены = [];
let координатыСтены = [];

function рисоватьСетку() {
  контекст.clearRect(0, 0, холст.width, холст.height);

  контекст.strokeStyle = 'lightgray';
  контекст.lineWidth = 0.5;

  const шагСетки = parseInt(полеШагаСетки.value) * 2;

  for (let x = 0; x < холст.width; x += шагСетки) {
    контекст.beginPath();
    контекст.moveTo(x, 0);
    контекст.lineTo(x, холст.height);
    контекст.stroke();
  }

  for (let y = 0; y < холст.height; y += шагСетки) {
    контекст.beginPath();
    контекст.moveTo(0, y);
    контекст.lineTo(холст.width, y);
    контекст.stroke();
  }
}
function началоРисованияСтены(событие) {
  if (!режимРисованияСтен) return;
  console.log("началоРисованияСтены");
  const прямоугольник = холст.getBoundingClientRect();
  началоX = Math.floor((событие.clientX - прямоугольник.left) / (parseInt(полеШагаСетки.value) * 2)) * (parseInt(полеШагаСетки.value) * 2);
  началоY = Math.floor((событие.clientY - прямоугольник.top) / (parseInt(полеШагаСетки.value) * 2)) * (parseInt(полеШагаСетки.value) * 2);

  координатыСтены = [{ x: началоX, y: началоY }];
  холст.removeEventListener('click', началоРисованияСтены); 
  холст.addEventListener('mousemove', предварительныйПросмотрСтены);
  холст.addEventListener('click', добавитьТочкуСтены);
  холст.addEventListener('dblclick', законитьРисованиеСтены);
}

let включёнПредварительныйПросмтор = false


function предварительныйПросмотрСтены(событие) {
  включёнПредварительныйПросмтор = true
  console.log("предварительныйПросмотрСтены");
  const прямоугольник = холст.getBoundingClientRect();
  const конецX = Math.floor((событие.clientX - прямоугольник.left) / (parseInt(полеШагаСетки.value) * 2)) * (parseInt(полеШагаСетки.value) * 2);
  const конецY = Math.floor((событие.clientY - прямоугольник.top) / (parseInt(полеШагаСетки.value) * 2)) * (parseInt(полеШагаСетки.value) * 2);

  контекст.clearRect(0, 0, холст.width, холст.height);
  рисоватьСетку();
  рисоватьСтены();

  контекст.strokeStyle = 'rgba(0, 0, 0)';
  контекст.lineWidth = parseInt(полеТолщиныСтены.value) * 2;
  контекст.beginPath();

  контекст.moveTo(координатыСтены[0].x, координатыСтены[0].y);


  for (let i = 1; i < координатыСтены.length; i++) {
    const толщина = parseInt(полеТолщиныСтены.value);
    const предыдущийX = координатыСтены[i - 1].x;
    const предыдущийY = координатыСтены[i - 1].y;
    const текущийX = координатыСтены[i].x;
    const текущийY = координатыСтены[i].y;

    контекст.lineTo(текущийX, текущийY);

  }
  контекст.moveTo(началоX, началоY);
  контекст.lineTo(конецX, конецY); 

  контекст.stroke();

  длинаСтены = Math.sqrt(Math.pow(конецX - началоX, 2) + Math.pow(конецY - началоY, 2)) / 2;
  полеДлиныСтены.textContent = длинаСтены.toFixed(2);
}


function рисоватьСтены() {

  console.log("рисоватьСтены");
  контекст.strokeStyle = 'black';
  контекст.lineWidth = parseInt(полеТолщиныСтены.value) * 2;
  // контекст.beginPath();
  комнаты.forEach(стена => {
   console.log(комнаты);
    контекст.beginPath();
    контекст.moveTo(стена[0].x, стена[0].y);

    for (let i = 1; i < стена.length; i++) {
      // const толщина = parseInt(полеТолщиныСтены.value);
      // const предыдущийX = стена[i - 1].x;
      // const предыдущийY = стена[i - 1].y;
      const текущийX = стена[i].x;
      const текущийY = стена[i].y;

      контекст.lineTo(текущийX, текущийY);
      // if (предыдущийX === текущийX) {
      //   // Вертикальная линия
      //   контекст.lineTo(текущийX - толщина, текущийY);
      // } else if (предыдущийY === текущийY) {
      //   // Горизонтальная линия
      //   контекст.lineTo(текущийX, текущийY - толщина);
      // }
    }

    контекст.closePath();
    // контекст.fillStyle = 'white';
    // контекст.fill();
    контекст.stroke();
  });
}


function добавитьТочкуСтены(событие) {


if (!включёнПредварительныйПросмтор ) {
  включёнПредварительныйПросмтор = true
  холст.addEventListener('mousemove', предварительныйПросмотрСтены);
  холст.addEventListener('dblclick', законитьРисованиеСтены);
}


  const прямоугольник = холст.getBoundingClientRect();
  const конецX = Math.floor((событие.clientX - прямоугольник.left) / (parseInt(полеШагаСетки.value) * 2)) * (parseInt(полеШагаСетки.value) * 2);
  const конецY = Math.floor((событие.clientY - прямоугольник.top) / (parseInt(полеШагаСетки.value) * 2)) * (parseInt(полеШагаСетки.value) * 2);

  координатыСтены.push({ x: конецX, y: конецY });
  началоX = конецX;
  началоY = конецY;

  контекст.clearRect(0, 0, холст.width, холст.height);
  рисоватьСетку();
  рисоватьСтены();

  контекст.strokeStyle = 'black';
  контекст.lineWidth = parseInt(полеТолщиныСтены.value) * 2;
  контекст.beginPath(); 
  контекст.moveTo(координатыСтены[0].x, координатыСтены[0].y);


  for (let i = 1; i < координатыСтены.length; i++) {
    const толщина = parseInt(полеТолщиныСтены.value);
    const предыдущийX = координатыСтены[i - 1].x;
    const предыдущийY = координатыСтены[i - 1].y;
    const текущийX = координатыСтены[i].x;
    const текущийY = координатыСтены[i].y;

    контекст.lineTo(текущийX, текущийY);

    // if (предыдущийX === текущийX) {
    //   // Вертикальная линия
    //   контекст.lineTo(текущийX - толщина, текущийY);
    // } else if (предыдущийY === текущийY) {
    //   // Горизонтальная линия
    //   контекст.lineTo(текущийX, текущийY - толщина);
    // }
  }

  контекст.stroke();

}

function законитьРисованиеСтены() {
  console.log("законитьРисованиеСтены");
  комнаты.push(координатыСтены);
  координатыСтены = [];

  холст.removeEventListener('mousemove', предварительныйПросмотрСтены);
  включёнПредварительныйПросмтор = false
  // холст.removeEventListener('click', добавитьТочкуСтены);
  // холст.removeEventListener('dblclick', законитьРисованиеСтены);

  вычислитьПлощадьКомнаты();

}







function вычислитьПлощадьКомнаты() {
  рисоватьСетку();
  // контекст.clearRect(0, 0, холст.width, холст.height);



 console.log("вычислитьПлощадьКомнаты");

  комнаты.forEach(стена => {
    console.log("стена", стена);

    контекст.beginPath();
    контекст.moveTo(стена[0].x, стена[0].y);

    for (let i = 1; i < стена.length; i++) {
      контекст.lineTo(стена[i].x, стена[i].y);
    }

    контекст.closePath();

    const площадь = Math.abs(получитьПлощадьМногоугольника(стена)) / 4; // Площадь в квадратных сантиметрах
    const текстПлощади = `${(площадь / 10000).toFixed(2)} м²`; // Площадь в квадратных метрах
  
    const центрX = стена.reduce((сумма, точка) => сумма + точка.x, 0) / стена.length;
    const центрY = стена.reduce((сумма, точка) => сумма + точка.y, 0) / стена.length;

    контекст.font = '16px Arial';
    контекст.fillStyle = 'black';
    контекст.textAlign = 'center';
    контекст.fillText(текстПлощади, центрX, центрY);
  });

  рисоватьСтены();


}

function получитьПлощадьМногоугольника(вершины) {
  let площадь = 0;

  for (let i = 0; i < вершины.length; i++) { const j = (i + 1) % вершины.length; площадь += вершины[i].x * вершины[j].y; площадь -= вершины[j].x * вершины[i].y; }

  return площадь / 2;
}



function очиститьСтены() {
  контекст.clearRect(0, 0, холст.width, холст.height);
  рисоватьСетку(); 
  стены = []; 
  полеДлиныСтены.textContent = '0';
}

кнопкаРисованияСтен.addEventListener('click', () => {
  режимРисованияСтен = !режимРисованияСтен; 
  режимРисованияДверей = false; 
  кнопкаРисованияСтен.textContent = режимРисованияСтен ? 'Закончить рисование стен' : 'Рисовать стену'; 
  кнопкаРисованияДверей.textContent = 'Рисовать дверь';

  if (режимРисованияСтен) {  
    холст.style.cursor = 'crosshair'; 
    холст.addEventListener('click', началоРисованияСтены); 

  } else { 
    холст.style.cursor = 'default'; 
     холст.removeEventListener('click', началоРисованияСтены); 
     холст.removeEventListener('mousemove', предварительныйПросмотрСтены);
     холст.removeEventListener('click', добавитьТочкуСтены);
     холст.removeEventListener('dblclick', законитьРисованиеСтены);
  }
});



кнопкаОчисткиСтен.addEventListener('click', очиститьСтены);

полеШагаСетки.addEventListener('change', () => { рисоватьСетку(); 
  if (режимРисованияСтен) { холст.addEventListener('click', началоРисованияСтены); } });
// ...

function началоРисованияДвери(событие) {
  if (!режимРисованияДверей) return;

  const прямоугольник = холст.getBoundingClientRect();
  const x = событие.clientX - прямоугольник.left;
  const y = событие.clientY - прямоугольник.top;

  const дверь = {
    x: x,
    y: y,
    ширина: parseInt(полеШириныДвери.value),
    высота: parseInt(полеТолщиныСтены.value)
  };

  const стена = найтиСтену(x, y);
  if (стена) {
    стена.двери.push(дверь);
    рисоватьДверь(дверь, стена);
  }
}

function найтиСтену(x, y) {
  for (let i = 0; i < стены.length; i++) {
    const стена = стены[i];
    if (x >= стена.x && x <= стена.x + стена.длина && y >= стена.y && y <= стена.y + стена.толщина) {
      return стена;
    }
  }
  return null;
}

function рисоватьДверь(дверь, стена) {
  контекст.fillStyle = 'white';
  if (стена.горизонтальная) {
    контекст.fillRect(дверь.x - дверь.ширина / 2, стена.y, дверь.ширина, дверь.высота);
  } else {
    контекст.fillRect(стена.x, дверь.y - дверь.ширина / 2, дверь.высота, дверь.ширина);
  }
}

function предварительныйПросмотрДвери(событие) {
  if (!режимРисованияДверей) return;

  const прямоугольник = холст.getBoundingClientRect();
  const x = событие.clientX - прямоугольник.left;
  const y = событие.clientY - прямоугольник.top;

  const дверь = {
    x: x,
    y: y,
    ширина: parseInt(полеШириныДвери.value),
    высота: parseInt(полеТолщиныСтены.value)
  };

  const стена = найтиСтену(x, y);
  if (стена) {
    контекст.clearRect(0, 0, холст.width, холст.height);
    рисоватьСетку();
    рисоватьСтены();
    рисоватьДверь(дверь, стена);
  }
}

// ...

кнопкаРисованияДверей.addEventListener('click', () => {
  режимРисованияДверей = !режимРисованияДверей;
  режимРисованияСтен = false;
  кнопкаРисованияДверей.textContent = режимРисованияДверей ? 'Закончить рисование дверей' : 'Рисовать дверь';
  кнопкаРисованияСтен.textContent = 'Рисовать стену';

  if (режимРисованияДверей) {
    холст.style.cursor = 'pointer';
    холст.addEventListener('click', началоРисованияДвери);
    холст.addEventListener('mousemove', предварительныйПросмотрДвери);
  } else {
    холст.style.cursor = 'default';
    холст.removeEventListener('click', началоРисованияДвери);
    холст.removeEventListener('mousemove', предварительныйПросмотрДвери); } });

рисоватьСетку();