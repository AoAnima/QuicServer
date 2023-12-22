const canvas = document.getElementById('canvas');
const ctx = canvas.getContext('2d');

const widthInput = document.getElementById('width');
const lengthInput = document.getElementById('length');

ctx.fillStyle = 'white';
ctx.fillRect(0, 0, canvas.width, canvas.height);

ctx.strokeStyle = 'black';
ctx.lineWidth = 2;
ctx.strokeRect(20, 20, widthInput.value, lengthInput.value);

widthInput.addEventListener('change', function () {
  ctx.strokeStyle = 'black';
  ctx.lineWidth = 2;
  ctx.strokeRect(20, 20, widthInput.value, lengthInput.value);
});

lengthInput.addEventListener('change', function () {
  ctx.strokeStyle = 'black';
  ctx.lineWidth = 2;
  ctx.strokeRect(20, 20, widthInput.value, lengthInput.value);
});

 const step = 10; // равный шаг между прямыми
 const totalLength = 100; // общая длина всех прямых
 const n = totalLength / step; // количество прямых
 const squareSize = 100; // размер квадрата

 for (let i = 0; i < n; i++) {
   ctx.beginPath();
   ctx.moveTo(i * step, 0);
   ctx.lineTo(i * step, squareSize);
   ctx.stroke();
 }


 const w = canvas.width;
 const h=canvas.height;

 const mouse = { x:0, y:0}; // координаты мыши
 let draw = false;

 // нажатие мыши
 canvas.addEventListener("mousedown", function(e){
   mouse.x = e.pageX - this.offsetLeft;
   mouse.y = e.pageY - this.offsetTop;
   draw = true;
   context.beginPath();
   context.moveTo(mouse.x, mouse.y);
 });

 // перемещение мыши
 canvas.addEventListener("mousemove", function(e){
   if(draw==true){
     mouse.x = e.pageX - this.offsetLeft;
     mouse.y = e.pageY - this.offsetTop;
     context.lineTo(mouse.x, mouse.y);
     context.stroke();
   }
 });

 // отпускаем мышь
 canvas.addEventListener("mouseup", function(e){
   mouse.x = e.pageX - this.offsetLeft;
   mouse.y = e.pageY - this.offsetTop;
   context.lineTo(mouse.x, mouse.y);
   context.stroke();
   context.closePath();
   draw = false;
 });