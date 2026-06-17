```javascript
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAchievementStore } from './achievement.js'

// Mock Date.now for consistent testing
const mockDateNow = 1234567890000
vi.spyOn(Date, 'now').mockReturnValue(mockDateNow)

describe('useAchievementStore', () => {
  let store

  beforeEach(() => {
    setActivePinia(createPinia())
    store = useAchievementStore()
  })

  describe('初始状态', () => {
    it('应该有10个初始成就', () => {
      expect(store.achievements).toHaveLength(10)
    })

    it('应该有正确的条件类型列表', () => {
      expect(store.conditionTypes).toHaveLength(10)
      expect(store.conditionTypes[0]).toEqual({
        value: 'quest_complete',
        label: '完成任务数',
        description: '完成指定数量的任务'
      })
    })

    it('第一个成就应该是"初来乍到"', () => {
      expect(store.achievements[0]).toEqual({
        id: 1,
        name: '初来乍到',
        code: 'ach_first_quest',
        description: '完成第一个任务',
        condition: { type: 'quest_complete', value: 1 },
        reward: { exp: 50, gold: 100 },
        icon: '⭐'
      })
    })
  })

  describe('addAchievement', () => {
    it('应该添加新成就并生成唯一ID', () => {
      const initialCount = store.achievements.length
      const newAchievement = {
        name: '新成就',
        code: 'ach_new',
        description: '这是一个新成就',
        condition: { type: 'test', value: 1 },
        reward: { exp: 100, gold: 200 },
        icon: '🎯'
      }

      store.addAchievement(newAchievement)

      expect(store.achievements).toHaveLength(initialCount + 1)
      const addedAchievement = store.achievements.find(a => a.code === 'ach_new')
      expect(addedAchievement).toBeDefined()
      expect(addedAchievement.id).toBe(mockDateNow)
      expect(addedAchievement.name).toBe('新成就')
    })

    it('应该为每个新成就生成不同的ID', () => {
      Date.now.mockReturnValueOnce(111).mockReturnValueOnce(222)
      
      store.addAchievement({ code: 'ach_1', name: '成就1' })
      store.addAchievement({ code: 'ach_2', name: '成就2' })

      expect(store.achievements.find(a => a.code === 'ach_1').id).toBe(111)
      expect(store.achievements.find(a => a.code === 'ach_2').id).toBe(222)
    })
  })

  describe('updateAchievement', () => {
    it('应该更新现有成就', () => {
      const achievementId = 1
      const updateData = {
        name: '更新后的成就',
        description: '更新后的描述'
      }

      store.updateAchievement(achievementId, updateData)

      const updatedAchievement = store.getAchievementById(achievementId)
      expect(updatedAchievement.name).toBe('更新后的成就')
      expect(updatedAchievement.description).toBe('更新后的描述')
      expect(updatedAchievement.code).toBe('ach_first_quest') // 保持未更新的字段
    })

    it('当ID不存在时不应更新任何成就', () => {
      const initialAchievements = [...store.achievements]
      const nonExistentId = 999

      store.updateAchievement(nonExistentId, { name: '不存在' })

      expect(store.achievements).toEqual(initialAchievements)
    })

    it('应该更新成就的奖励', () => {
      const achievementId = 2
      const newReward = { exp: 500, gold: 1000 }

      store.updateAchievement(achievementId, { reward: newReward })

      const updatedAchievement = store.getAchievementById(achievementId)
      expect(updatedAchievement.reward).toEqual(newReward)
    })
  })

  describe('deleteAchievement', () => {
    it('应该删除指定ID的成就', () => {
      const initialCount = store.achievements.length
      const achievementIdToDelete = 3

      store.deleteAchievement(achievementIdToDelete)

      expect(store.achievements).toHaveLength(initialCount - 1)
      expect(store.getAchievementById(achievementIdToDelete)).toBeUndefined()
    })

    it('当ID不存在时不应删除任何成就', () => {
      const initialCount = store.achievements.length
      const nonExistentId = 999

      store.deleteAchievement(nonExistentId)

      expect(store.achievements).toHaveLength(initialCount)
    })

    it('删除后其他成就应保持不变', () => {
      const achievementToDelete = store.achievements.find(a => a.id === 1)
      const otherAchievements = store.achievements.filter(a => a.id !== 1)

      store.deleteAchievement(1)

      otherAchievements.forEach(achievement => {
        expect(store.getAchievementById(achievement.id)).toEqual(achievement)
      })
    })
  })

  describe('getAchievementById', () => {
    it('应该返回指定ID的成就', () => {
      const achievement = store.getAchievementById(1)
      expect(achievement).toBeDefined()
      expect(achievement.id).toBe(1)
      expect(achievement.name).toBe('初来乍到')
    })

    it('当ID不存在时应返回undefined', () => {
      const nonExistentId = 999
      const result = store.getAchievementById(nonExistentId)
      expect(result).toBeUndefined()
    })

    it('应该返回成就的完整对象', () => {
      const achievement = store.getAchievementById(2)
      expect(achievement).toEqual({
        id: 2,
        name: '村庄之友',
        code: 'ach_talk_all_npcs',
        description: '与所有NPC对话过',
        condition: { type: 'talk_all_npcs', value: 1 },
        reward: { exp: 200, gold: 300 },
        icon: '👥'
      })
    })
  })
})
```