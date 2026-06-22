```javascript
import { defineStore } from 'pinia'

export const useAchievementStore = defineStore('achievement', {
  state: () => ({
    achievements: []
  }),
  getters: {
    getAchievementCount: (state) => state.achievements.length,
    getAchievementByCode: (state) => {
      return (code) => {
        return state.achievements.find(a => String(a.code) === String(code)) || null
      }
    }
  },
  actions: {
    addAchievement(achievement) {
      const existing = this.achievements.find(a => String(a.code) === String(achievement.code))
      if (!existing) {
        this.achievements.push(achievement)
      }
    },
    removeAchievement(code) {
      const index = this.achievements.findIndex(achievement => String(achievement.code) === String(code))
      if (index !== -1) {
        this.achievements.splice(index, 1)
      }
    },
    resetAchievements() {
      this.achievements = []
    },
    addMultipleAchievements(achievements) {
      const uniqueAchievements = []
      const existingCodes = new Set(this.achievements.map(a => String(a.code)))
      
      achievements.forEach(achievement => {
        if (!existingCodes.has(String(achievement.code))) {
          uniqueAchievements.push(achievement)
          existingCodes.add(String(achievement.code))
        }
      })
      
      this.achievements.push(...uniqueAchievements)
    }
  }
})
```