export class MiniMap {
    constructor(scene) {
        this.scene = scene;
        this.container = null;
        this.mapSize = 150;
        this.scale = 0.1;
    }

    create() {
        if (this.container) {
            this.container.destroy();
        }

        const width = this.scene.cameras.main.width;
        const height = this.scene.cameras.main.height;
        const x = width - this.mapSize - 10;
        const y = 10;

        this.container = this.scene.add.container(x, y).setDepth(250).setScrollFactor(0);

        const bg = this.scene.add.rectangle(this.mapSize / 2, this.mapSize / 2, this.mapSize, this.mapSize, 0x000000, 0.6).setStrokeStyle(1, 0x888888);
        this.container.add(bg);

        this.sceneData = this.scene.gameData.scenes.find(s => s.code === this.scene.currentSceneCode);
        if (this.sceneData) {
            this.scale = this.mapSize / Math.max(this.sceneData.width, this.sceneData.height);
            this.drawSceneBounds();
        }

        this.playerDot = this.scene.add.circle(0, 0, 4, 0x00ff00);
        this.container.add(this.playerDot);

        this.npcDots = [];
        this.portalDots = [];

        this.update();
    }

    drawSceneBounds() {
        if (!this.sceneData) return;

        const w = this.sceneData.width * this.scale;
        const h = this.sceneData.height * this.scale;

        const border = this.scene.add.rectangle(w / 2, h / 2, w, h).setStrokeStyle(1, 0x555555).setFillStyle();
        this.container.add(border);
    }

    update() {
        if (!this.container || !this.scene.player) return;

        this.npcDots.forEach(dot => dot.destroy());
        this.portalDots.forEach(dot => dot.destroy());
        this.npcDots = [];
        this.portalDots = [];

        const playerX = this.scene.player.x * this.scale;
        const playerY = this.scene.player.y * this.scale;
        this.playerDot.setPosition(playerX, playerY);

        this.scene.npcSprites.forEach((npcSprite) => {
            const dot = this.scene.add.circle(
                npcSprite.x * this.scale,
                npcSprite.y * this.scale,
                3, 0x3498db
            );
            this.container.add(dot);
            this.npcDots.push(dot);
        });

        this.scene.portalSprites.forEach((portalSprite) => {
            const dot = this.scene.add.circle(
                portalSprite.x * this.scale,
                portalSprite.y * this.scale,
                3, 0x9b59b6
            );
            this.container.add(dot);
            this.portalDots.push(dot);
        });
    }

    destroy() {
        if (this.container) {
            this.container.destroy();
            this.container = null;
        }
    }
}