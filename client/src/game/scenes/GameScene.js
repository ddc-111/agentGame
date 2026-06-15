import Phaser from 'phaser';

export class GameScene extends Phaser.Scene {
    constructor() {
        super({ key: 'GameScene' });
    }

    create() {
        this.add.text(
            this.cameras.main.centerX,
            this.cameras.main.centerY,
            '古风RPG - Agent游戏',
            {
                fontSize: '32px',
                fontFamily: 'Arial',
                color: '#ffffff'
            }
        ).setOrigin(0.5);
    }

    update() {
    }
}
