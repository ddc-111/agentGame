import Phaser from 'phaser';

export class BootScene extends Phaser.Scene {
    constructor() {
        super({ key: 'BootScene' });
    }

    preload() {
        // Create loading bar
        const width = this.cameras.main.width;
        const height = this.cameras.main.height;
        
        const progressBar = this.add.graphics();
        const progressBox = this.add.graphics();
        progressBox.fillStyle(0x222222, 0.8);
        progressBox.fillRect(width / 2 - 160, height / 2 - 25, 320, 50);
        
        const loadingText = this.add.text(width / 2, height / 2 - 50, '加载中...', {
            font: '20px Microsoft YaHei',
            fill: '#ffffff'
        }).setOrigin(0.5, 0.5);

        this.load.on('progress', (value) => {
            progressBar.clear();
            progressBar.fillStyle(0x4a7c59, 1);
            progressBar.fillRect(width / 2 - 150, height / 2 - 15, 300 * value, 30);
        });

        this.load.on('complete', () => {
            progressBar.destroy();
            progressBox.destroy();
            loadingText.destroy();
        });

        // Generate textures programmatically
        this.createTextures();
    }

    createTextures() {
        // Player texture - simple character
        const playerCanvas = this.textures.createCanvas('player', 32, 48);
        const pCtx = playerCanvas.getContext();
        // Body
        pCtx.fillStyle = '#3498db';
        pCtx.fillRect(8, 16, 16, 24);
        // Head
        pCtx.fillStyle = '#f5cba7';
        pCtx.fillRect(10, 2, 12, 14);
        // Eyes
        pCtx.fillStyle = '#2c3e50';
        pCtx.fillRect(13, 7, 2, 2);
        pCtx.fillRect(18, 7, 2, 2);
        // Hair
        pCtx.fillStyle = '#2c3e50';
        pCtx.fillRect(10, 2, 12, 4);
        // Feet
        pCtx.fillStyle = '#8b4513';
        pCtx.fillRect(10, 40, 5, 6);
        pCtx.fillRect(18, 40, 5, 6);
        playerCanvas.refresh();

        // NPC textures
        this.createNPCTexture('npc_chief', '#8b0000', '#f5cba7', '#c0c0c0');
        this.createNPCTexture('npc_merchant', '#b8860b', '#f5cba7', '#2c3e50');
        this.createNPCTexture('npc_tea', '#ff69b4', '#f5cba7', '#2c3e50');
        this.createNPCTexture('npc_blacksmith', '#696969', '#deb887', '#000000');
        this.createNPCTexture('npc_hunter', '#556b2f', '#deb887', '#2c3e50');
        this.createNPCTexture('npc_kid', '#ff6347', '#f5cba7', '#ffa500');

        // Tile textures
        this.createTileTexture('tile_grass', '#4a7c59');
        this.createTileTexture('tile_dirt', '#8b7355');
        this.createTileTexture('tile_stone', '#808080');
        this.createTileTexture('tile_wood', '#6b4226');
        this.createTileTexture('tile_water', '#4a90d9');

        // Portal texture
        const portalCanvas = this.textures.createCanvas('portal', 48, 48);
        const ptCtx = portalCanvas.getContext();
        ptCtx.fillStyle = '#9b59b6';
        ptCtx.globalAlpha = 0.6;
        ptCtx.beginPath();
        ptCtx.arc(24, 24, 20, 0, Math.PI * 2);
        ptCtx.fill();
        ptCtx.globalAlpha = 1;
        ptCtx.strokeStyle = '#8e44ad';
        ptCtx.lineWidth = 2;
        ptCtx.stroke();
        portalCanvas.refresh();

        // Tree texture
        const treeCanvas = this.textures.createCanvas('tree', 48, 64);
        const tCtx = treeCanvas.getContext();
        // Trunk
        tCtx.fillStyle = '#8b4513';
        tCtx.fillRect(20, 40, 8, 24);
        // Leaves
        tCtx.fillStyle = '#228b22';
        tCtx.beginPath();
        tCtx.moveTo(24, 0);
        tCtx.lineTo(48, 40);
        tCtx.lineTo(0, 40);
        tCtx.closePath();
        tCtx.fill();
        treeCanvas.refresh();

        // House texture
        const houseCanvas = this.textures.createCanvas('house', 96, 80);
        const hCtx = houseCanvas.getContext();
        // Wall
        hCtx.fillStyle = '#deb887';
        hCtx.fillRect(8, 30, 80, 50);
        // Roof
        hCtx.fillStyle = '#8b0000';
        hCtx.beginPath();
        hCtx.moveTo(0, 30);
        hCtx.lineTo(48, 0);
        hCtx.lineTo(96, 30);
        hCtx.closePath();
        hCtx.fill();
        // Door
        hCtx.fillStyle = '#654321';
        hCtx.fillRect(38, 50, 20, 30);
        // Window
        hCtx.fillStyle = '#87ceeb';
        hCtx.fillRect(15, 42, 15, 15);
        hCtx.fillRect(66, 42, 15, 15);
        houseCanvas.refresh();

        // NPC interaction indicator
        const indicatorCanvas = this.textures.createCanvas('indicator', 24, 24);
        const iCtx = indicatorCanvas.getContext();
        iCtx.fillStyle = '#f1c40f';
        iCtx.beginPath();
        iCtx.moveTo(12, 0);
        iCtx.lineTo(24, 24);
        iCtx.lineTo(0, 24);
        iCtx.closePath();
        iCtx.fill();
        iCtx.fillStyle = '#2c3e50';
        iCtx.font = '16px Arial';
        iCtx.textAlign = 'center';
        iCtx.fillText('!', 12, 20);
        indicatorCanvas.refresh();

        // Enemy textures
        this.createEnemyTexture('enemy_wolf', '#7f8c8d', '#ecf0f1');
        this.createEnemyTexture('enemy_alpha_wolf', '#5d6d7e', '#e74c3c');
    }

    createEnemyTexture(key, bodyColor, eyeColor) {
        const canvas = this.textures.createCanvas(key, 48, 48);
        const ctx = canvas.getContext();
        // Body
        ctx.fillStyle = bodyColor;
        ctx.beginPath();
        ctx.ellipse(24, 28, 18, 12, 0, 0, Math.PI * 2);
        ctx.fill();
        // Head
        ctx.fillStyle = bodyColor;
        ctx.beginPath();
        ctx.arc(36, 18, 10, 0, Math.PI * 2);
        ctx.fill();
        // Ears
        ctx.fillStyle = bodyColor;
        ctx.beginPath();
        ctx.moveTo(30, 10);
        ctx.lineTo(28, 2);
        ctx.lineTo(34, 8);
        ctx.closePath();
        ctx.fill();
        ctx.beginPath();
        ctx.moveTo(40, 8);
        ctx.lineTo(42, 2);
        ctx.lineTo(44, 10);
        ctx.closePath();
        ctx.fill();
        // Eyes
        ctx.fillStyle = eyeColor;
        ctx.beginPath();
        ctx.arc(38, 16, 2, 0, Math.PI * 2);
        ctx.fill();
        // Nose
        ctx.fillStyle = '#2c3e50';
        ctx.beginPath();
        ctx.arc(44, 18, 2, 0, Math.PI * 2);
        ctx.fill();
        // Legs
        ctx.fillStyle = bodyColor;
        ctx.fillRect(12, 36, 4, 10);
        ctx.fillRect(20, 36, 4, 10);
        ctx.fillRect(28, 36, 4, 10);
        ctx.fillRect(36, 36, 4, 10);
        canvas.refresh();
    }

    createNPCTexture(key, bodyColor, skinColor, hairColor) {
        const canvas = this.textures.createCanvas(key, 32, 48);
        const ctx = canvas.getContext();
        // Body
        ctx.fillStyle = bodyColor;
        ctx.fillRect(8, 16, 16, 24);
        // Head
        ctx.fillStyle = skinColor;
        ctx.fillRect(10, 2, 12, 14);
        // Eyes
        ctx.fillStyle = '#2c3e50';
        ctx.fillRect(13, 7, 2, 2);
        ctx.fillRect(18, 7, 2, 2);
        // Hair
        ctx.fillStyle = hairColor;
        ctx.fillRect(10, 2, 12, 4);
        // Feet
        ctx.fillStyle = '#654321';
        ctx.fillRect(10, 40, 5, 6);
        ctx.fillRect(18, 40, 5, 6);
        canvas.refresh();
    }

    createTileTexture(key, color) {
        const canvas = this.textures.createCanvas(key, 48, 48);
        const ctx = canvas.getContext();
        ctx.fillStyle = color;
        ctx.fillRect(0, 0, 48, 48);
        // Grid lines
        ctx.strokeStyle = 'rgba(0,0,0,0.1)';
        ctx.strokeRect(0, 0, 48, 48);
        canvas.refresh();
    }

    create() {
        // Initialize game state
        this.registry.set('gameData', null);
        this.registry.set('playerData', null);
        this.registry.set('currentScene', null);
        
        this.scene.start('GameScene');
    }
}
