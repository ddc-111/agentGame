import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useDemoStore } from '@/stores/demo'

describe('useDemoStore', () => {
  let store
  
  beforeEach(() => {
    setActivePinia(createPinia())
    store = useDemoStore()
  })
  
  afterEach(() => {
    vi.clearAllMocks()
  })

  describe('initial state', () => {
    it('should have demos array', () => {
      expect(Array.isArray(store.demos)).toBe(true)
    })

    it('should have default reactive values', () => {
      expect(store.activeDemo).toBeNull()
      expect(store.currentStep).toBe(0)
      expect(store.isPlaying).toBe(false)
      expect(store.playbackSpeed).toBe(1)
    })
    
    it('should have demos populated', () => {
      expect(store.demos.length).toBe(3)
      expect(store.demos[0].id).toBe('beginner_village')
    })
  })

  describe('selectDemo', () => {
    it('should select a demo by id', () => {
      store.selectDemo('beginner_village')
      expect(store.activeDemo).toBeDefined()
      expect(store.activeDemo.id).toBe('beginner_village')
      expect(store.currentStep).toBe(0)
    })

    it('should set activeDemo to null for invalid id', () => {
      store.selectDemo('invalid_id')
      expect(store.activeDemo).toBeNull()
    })

    it('should reset currentStep to 0 when selecting demo', () => {
      store.currentStep = 5
      store.selectDemo('beginner_village')
      expect(store.currentStep).toBe(0)
    })

    it('should handle empty demos array', () => {
      store.demos = []
      store.selectDemo('beginner_village')
      expect(store.activeDemo).toBeNull()
    })
  })

  describe('nextStep', () => {
    it('should increment current step when valid steps exist', () => {
      store.selectDemo('beginner_village')
      store.nextStep()
      expect(store.currentStep).toBe(1)
    })

    it('should not exceed max steps', () => {
      store.selectDemo('beginner_village')
      store.currentStep = 13 // Last step (14 steps, 0-indexed)
      store.nextStep()
      expect(store.currentStep).toBe(13)
    })

    it('should not increment when no active demo', () => {
      store.currentStep = 0
      store.nextStep()
      expect(store.currentStep).toBe(0)
    })
  })

  describe('prevStep', () => {
    it('should decrement current step when not at first step', () => {
      store.selectDemo('beginner_village')
      store.currentStep = 2
      store.prevStep()
      expect(store.currentStep).toBe(1)
    })

    it('should not go below 0', () => {
      store.selectDemo('beginner_village')
      store.currentStep = 0
      store.prevStep()
      expect(store.currentStep).toBe(0)
    })
    
    it('should decrement even when no active demo', () => {
      store.currentStep = 5
      store.prevStep()
      expect(store.currentStep).toBe(4)
    })
  })

  describe('togglePlay', () => {
    it('should toggle isPlaying state', () => {
      store.togglePlay()
      expect(store.isPlaying).toBe(true)
      
      store.togglePlay()
      expect(store.isPlaying).toBe(false)
    })
    
    it('should start from beginning when no active demo', () => {
      store.togglePlay()
      expect(store.isPlaying).toBe(true)
      expect(store.currentStep).toBe(0)
    })
  })

  describe('stop', () => {
    it('should stop playback and reset current step', () => {
      store.selectDemo('beginner_village')
      store.isPlaying = true
      store.currentStep = 5

      store.stop()
      expect(store.isPlaying).toBe(false)
      expect(store.currentStep).toBe(0)
    })

    it('should work when already stopped', () => {
      store.isPlaying = false
      store.currentStep = 3

      store.stop()
      expect(store.isPlaying).toBe(false)
      expect(store.currentStep).toBe(0)
    })
  })

  describe('setSpeed', () => {
    it('should set playback speed', () => {
      store.setSpeed(2)
      expect(store.playbackSpeed).toBe(2)
    })
    
    it('should accept different speed values', () => {
      store.setSpeed(0.5)
      expect(store.playbackSpeed).toBe(0.5)
      
      store.setSpeed(4)
      expect(store.playbackSpeed).toBe(4)
    })
    
    it('should handle edge cases', () => {
      store.setSpeed(0)
      expect(store.playbackSpeed).toBe(0)
      
      store.setSpeed(10)
      expect(store.playbackSpeed).toBe(10)
    })
  })

  describe('integration tests', () => {
    it('should handle full workflow: select, play, step, stop', () => {
      store.selectDemo('beginner_village')
      expect(store.activeDemo.id).toBe('beginner_village')
      expect(store.currentStep).toBe(0)

      store.togglePlay()
      expect(store.isPlaying).toBe(true)

      store.nextStep()
      expect(store.currentStep).toBe(1)

      store.nextStep()
      expect(store.currentStep).toBe(2)

      store.prevStep()
      expect(store.currentStep).toBe(1)

      store.stop()
      expect(store.isPlaying).toBe(false)
      expect(store.currentStep).toBe(0)
    })

    it('should maintain demo selection after stopping', () => {
      store.selectDemo('combat_demo')
      store.togglePlay()
      store.nextStep()

      store.stop()
      expect(store.activeDemo.id).toBe('combat_demo')
    })
  })
})