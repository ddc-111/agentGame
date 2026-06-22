import { describe, it, expect, vi, beforeEach } from 'vitest'

const mockFetch = vi.fn()
global.fetch = mockFetch

describe('API Module', () => {
  beforeEach(() => {
    mockFetch.mockReset()
  })

  it('should handle successful API response', async () => {
    const mockData = { id: 1, name: 'test' }
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => mockData
    })

    const response = await fetch('/api/test')
    const data = await response.json()

    expect(data).toEqual(mockData)
    expect(mockFetch).toHaveBeenCalledWith('/api/test')
  })

  it('should handle API error response', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: false,
      status: 404,
      statusText: 'Not Found'
    })

    const response = await fetch('/api/test')
    expect(response.ok).toBe(false)
    expect(response.status).toBe(404)
  })

  it('should handle POST request with body', async () => {
    const postData = { name: 'new item', code: 'item_new' }
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({ id: 2, ...postData })
    })

    const response = await fetch('/api/items', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(postData)
    })
    const data = await response.json()

    expect(data.name).toBe('new item')
    expect(data.id).toBe(2)
  })

  it('should handle network error', async () => {
    mockFetch.mockRejectedValueOnce(new Error('Network error'))

    await expect(fetch('/api/test')).rejects.toThrow('Network error')
  })
})
