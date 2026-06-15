import Phaser from 'phaser';

export class BootScene extends Phaser.Scene {
    constructor() {
        super({ key: 'BootScene' });
    }

    preload() {
        this.load.image('logo', 'assets/logo.png');
    }

    create() {
        this.scene.start('GameScene');
    }
}
