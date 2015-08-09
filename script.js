var canvas = document.getElementById("mainCanvas");
var context = canvas.getContext("2d");
var serversocket = new WebSocket("ws://localhost:4000/ws");

var keys = [];

var width = 500, height = 400, speed = 1, score = 0;

var player = {
	x: 40,
	y: 40,
	width: 20,
	height: 20
};

serversocket.onopen = function() {
    serversocket.send("init");
}
serversocket.onmessage = function(e){
	var server_message = e.data;
	var res = server_message.split(",");
	player.x = parseInt(res[0]);
	player.y = parseInt(res[1]);
}

var player1 = {
	x: 40,
	y: 340,
	width: 20,
	height: 20
};

var player2 = {
	x: 560,
	y: 40,
	width: 20,
	height: 20
};

var player3 = {
	x: 560,
	y: 460,
	width: 20,
	height: 20
};

var cube = {
	x: Math.random() * (width - 20), //Note: Math.random() gives value b/w 0-1
	y: Math.random() * (height - 20),
	width: 20,
	height: 20
};


window.addEventListener("keydown", function(e){
	keys[e.keyCode] = true;
}, false);

window.addEventListener("keyup", function(e){
	delete keys[e.keyCode];
}, false);

function update(){
	if(keys[87]){
		player.y-=speed; //up
	}
	if(keys[83]){
		player.y+=speed; //down
	}
	if(keys[65]){
		player.x-=speed; //left
	}
	if(keys[68]){
		player.x+=speed; //right
	}
	
	if(player.x < 0){
		player.x = 0; //if on left edge
	}
	if(player.y < 0){
		player.y = 0; //if on top edge
	}
	if(player.x >= width - player.width){
		player.x = width - player.width; //if on right edge
	}
	if(player.y >= height - player.height){
		player.y = height - player.height; //if on bottom edge
	}
	
	var toSend = "0 " + player.x.toString() + " " + player.y.toString();
	serversocket.send(toSend);
	
	if(collision(player, cube)) process_collision();
	
	serversocket.send("2"); //Get server update
	serversocket.onmessage = function(e){
		var server_message = e.data;
		var res = server_message.split(",");
		
		var player_specs1 = res[0].split(" ");
		player1.x = parseInt(player_specs1[1]);
		player1.y = parseInt(player_specs1[2]);
		
		var player_specs2 = res[1].split(" ");
		player2.x = parseInt(player_specs2[1]);
		player2.y = parseInt(player_specs2[2]);
		
		var player_specs3 = res[2].split(" ");
		player3.x = parseInt(player_specs3[1]);
		player3.y = parseInt(player_specs3[2]);
		
		var cube_specs = res[3].split(" ");
		cube.x = parseInt(cube_specs[1])
		cube.y = parseInt(cube_specs[2])
	}
}

function render(){
	context.clearRect(0,0, width, height);
	
	context.fillStyle = "blue";
	context.fillRect(player.x, player.y, player.width, player.height);
	
	context.fillStyle = "red";
	context.fillRect(player1.x, player1.y, player1.width, player1.height);
	context.fillRect(player2.x, player2.y, player2.width, player2.height);
	context.fillRect(player3.x, player3.y, player3.width, player3.height);
	
	context.fillStyle = "green";
	context.fillRect(cube.x, cube.y, cube.width, cube.height);
	
	context.font = "bold 30px helvetica";
	context.fillText(score, 10, 30);
}

function process_collision(){
	score++;
	serversocket.send("1");
	//cube.x = Math.random() * (width - 20);
	//cube.y = Math.random() * (height - 20);
}

function collision(first, second){
	return !(first.x > second.x + second.width || 
		first.x + first.width < second.x ||
		first.y > second.y + second.height ||
		first.y + first.height < second.y);
}

function game(){
	update();
	render();
}

setInterval(function(){
	game();
}, 1000/45) //1000/30 = 30 fps