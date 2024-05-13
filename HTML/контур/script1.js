const холст = document.getElementById('холст');
const контекст = холст.getContext('2d');
const полеШагаСетки = document.getElementById('шагСетки');
const полеТолщиныСтены = document.getElementById('толщинаСтен');
const кнопкаРисованияСтен = document.getElementById('рисоватьСтену');
const рисоватьДверь = document.getElementById('рисоватьДверь');
const полеДлиныСтены = document.getElementById('длинаСтены');
const полеШириныДвери = document.getElementById('ширинаДвери');

let режимРисованияСтен = false;
let режимРисованияДверей = false;
let началоX, началоY;
let длинаСтены = 0;

let стены = [];
let координатыСтены = [];

console.log('холст:', холст);
console.log('контекст:', контекст);
console.log('полеШагаСетки:', полеШагаСетки);



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


  


 холст.addEventListener('mousemove', (e) => {
  if (режимРисованияСтен) {
    const endX = e.offsetX;
    const endY = e.offsetY;
    контекст.clearRect(0, 0, холст.width, холст.height);
    рисоватьСетку();
    рисоватьСтены();
    контекст.strokeStyle = 'black';
    контекст.lineWidth = 2;
    контекст.beginPath();
    контекст.moveTo(началоX, началоY);
    контекст.lineTo(endX, endY);
    контекст.stroke();
    длинаСтены = Math.sqrt(Math.pow(endX - началоX, 2) + Math.pow(endY - началоY, 2));
    document.getElementById('длинаСтены').innerHTML = `${длинаСтены.toFixed(2)} см`;
  }
});

холст.addEventListener('click', (e) => {
  console.log('первый click');
  if (режимРисованияСтен) {
    началоX = e.offsetX;
    началоY = e.offsetY;
    координатыСтены = [];
    координатыСтены.push([началоX, началоY]);
  }
});

 холст.addEventListener('click', () => {
  console.log('второй click');
  if (режимРисованияСтен) {
    координатыСтены.push([началоX, началоY]);
    стены.push(координатыСтены);
    координатыСтены = [];
  }
});

 холст.addEventListener('dblclick', () => {
  console.log('dblclick');

  if (режимРисованияСтен) {
    стены.push(координатыСтены);
    координатыСтены = [];
  }
});

function рисоватьСтены() {
  контекст.strokeStyle = 'black';
  контекст.lineWidth = 2;
  for (const стена of стены) {
    контекст.beginPath();
    for (let i = 0; i < стена.length - 1; i++) {
      контекст.moveTo(стена[i][0], стена[i][1]);
      контекст.lineTo(стена[i + 1][0], стена[i + 1][1]);
    }
    контекст.stroke();
  }
}

кнопкаРисованияСтен.addEventListener('click', () => {
  режимРисованияСтен = true;
  document.body.style.cursor = 'crosshair';
});

// document.addEventListener('click', () => {
//   if (режимРисованияСтен) {
//     режимРисованияСтен = false;
//     document.body.style.cursor = 'default';
//   }
 
// });


  // Вызываем функцию рисования сетки при изменении шага сетки
  полеШагаСетки.addEventListener('input', рисоватьСетку);
  
  // Рисуем сетку при загрузке страницы
  рисоватьСетку();


function началоРисованияДвери(событие) {
  if (!режимРисованияДверей) return;

  const прямоугольник = холст.getBoundingClientRect(); const x = Math.floor((событие.clientX - прямоугольник.left) / (parseInt(полеШагаСетки.value) * 2)) * (parseInt(полеШагаСетки.value) * 2); const y = Math.floor((событие.clientY - прямоугольник.top) / (parseInt(полеШагаСетки.value) * 2)) * (parseInt(полеШагаСетки.value) * 2);

  const ширинаДвери = parseInt(полеШириныДвери.value) * 2;

  стены.forEach(стена => {
    for (let i = 1; i < стена.length; i++) {
      const началоX = стена[i - 1].x; const началоY = стена[i - 1].y; const конецX = стена[i].x; const конецY = стена[i].y;

      if (началоX === конецX && Math.abs(y - началоY) <= ширинаДвери / 2) {
        // Вертикальная стена
        контекст.fillStyle = 'gray';
        контекст.fillRect(началоX - ширинаДвери / 2, y - ширинаДвери / 2, ширинаДвери, ширинаДвери);
      } else if (началоY === конецY && Math.abs(x - началоX) <= ширинаДвери / 2) {
        // Горизонтальная стена
        контекст.fillStyle = 'gray';
        контекст.fillRect(x - ширинаДвери / 2, началоY - ширинаДвери / 2, ширинаДвери, ширинаДвери);
      }
    }
  });
}
