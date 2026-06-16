package network

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/database/models"
	"github.com/ddc-111/agentGame/server/internal/database/repository"
)

func (s *Server) handleExport(c *gin.Context) {
	data := make(map[string]interface{})

	scenes, _ := s.repo.GetScenes()
	data["scenes"] = scenes

	npcs, _ := s.repo.GetNPCs()
	data["npcs"] = npcs

	agents, _ := s.repo.GetAgents()
	data["agents"] = agents

	shops, _ := s.repo.GetShops()
	data["shops"] = shops

	items, _ := s.repo.GetItems()
	data["items"] = items

	tasks, _ := s.repo.GetTasks()
	data["tasks"] = tasks

	flows, _ := s.repo.GetFlows()
	data["flows"] = flows

	templates, _ := s.repo.GetTemplates()
	data["prompts"] = templates

	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (s *Server) handleImport(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	if inner, ok := data["data"].(map[string]interface{}); ok {
		data = inner
	}

	imported := make(map[string]int)

	err := s.repo.Transaction(func(repo *repository.Repository) error {
		if scenesRaw, ok := data["scenes"].([]interface{}); ok {
			for _, sr := range scenesRaw {
				sceneData, _ := jsonMarshal(sr)
				var scene models.Scene
				if jsonUnmarshal(sceneData, &scene) == nil {
					if existing, _ := repo.GetSceneByCode(scene.Code); existing != nil && existing.ID > 0 {
						scene.ID = existing.ID
						if err := repo.UpdateScene(&scene); err != nil {
							return err
						}
					} else {
						scene.ID = 0
						if err := repo.CreateScene(&scene); err != nil {
							return err
						}
					}
					imported["scenes"]++
				}
			}
		}

		if npcsRaw, ok := data["npcs"].([]interface{}); ok {
			for _, nr := range npcsRaw {
				npcData, _ := jsonMarshal(nr)
				var npc models.NPC
				if jsonUnmarshal(npcData, &npc) == nil {
					if existing, _ := repo.GetNPCByCode(npc.Code); existing != nil && existing.ID > 0 {
						npc.ID = existing.ID
						if err := repo.UpdateNPC(&npc); err != nil {
							return err
						}
					} else {
						npc.ID = 0
						if err := repo.CreateNPC(&npc); err != nil {
							return err
						}
					}
					imported["npcs"]++
				}
			}
		}

		if agentsRaw, ok := data["agents"].([]interface{}); ok {
			for _, ar := range agentsRaw {
				agentData, _ := jsonMarshal(ar)
				var agent models.Agent
				if jsonUnmarshal(agentData, &agent) == nil {
					if existing, _ := repo.GetAgentByCode(agent.Code); existing != nil && existing.ID > 0 {
						agent.ID = existing.ID
						if err := repo.UpdateAgent(&agent); err != nil {
							return err
						}
					} else {
						agent.ID = 0
						if err := repo.CreateAgent(&agent); err != nil {
							return err
						}
					}
					imported["agents"]++
				}
			}
		}

		if shopsRaw, ok := data["shops"].([]interface{}); ok {
			for _, sr := range shopsRaw {
				shopData, _ := jsonMarshal(sr)
				var shop models.Shop
				if jsonUnmarshal(shopData, &shop) == nil {
					if existing, _ := repo.GetShopByCode(shop.Code); existing != nil && existing.ID > 0 {
						shop.ID = existing.ID
						if err := repo.UpdateShop(&shop); err != nil {
							return err
						}
					} else {
						shop.ID = 0
						if err := repo.CreateShop(&shop); err != nil {
							return err
						}
					}
					imported["shops"]++
				}
			}
		}

		if itemsRaw, ok := data["items"].([]interface{}); ok {
			for _, ir := range itemsRaw {
				itemData, _ := jsonMarshal(ir)
				var item models.Item
				if jsonUnmarshal(itemData, &item) == nil {
					if existing, _ := repo.GetItemByCode(item.Code); existing != nil && existing.ID > 0 {
						item.ID = existing.ID
						if err := repo.UpdateItem(&item); err != nil {
							return err
						}
					} else {
						item.ID = 0
						if err := repo.CreateItem(&item); err != nil {
							return err
						}
					}
					imported["items"]++
				}
			}
		}

		if tasksRaw, ok := data["tasks"].([]interface{}); ok {
			for _, tr := range tasksRaw {
				taskData, _ := jsonMarshal(tr)
				var task models.Task
				if jsonUnmarshal(taskData, &task) == nil {
					if existing, _ := repo.GetTaskByCode(task.Code); existing != nil && existing.ID > 0 {
						task.ID = existing.ID
						if err := repo.UpdateTask(&task); err != nil {
							return err
						}
					} else {
						task.ID = 0
						if err := repo.CreateTask(&task); err != nil {
							return err
						}
					}
					imported["tasks"]++
				}
			}
		}

		if flowsRaw, ok := data["flows"].([]interface{}); ok {
			for _, fr := range flowsRaw {
				flowData, _ := jsonMarshal(fr)
				var flow models.Flow
				if jsonUnmarshal(flowData, &flow) == nil {
					if existing, _ := repo.GetFlowByCode(flow.Code); existing != nil && existing.ID > 0 {
						flow.ID = existing.ID
						if err := repo.UpdateFlow(&flow); err != nil {
							return err
						}
					} else {
						flow.ID = 0
						if err := repo.CreateFlow(&flow); err != nil {
							return err
						}
					}
					imported["flows"]++
				}
			}
		}

		if templatesRaw, ok := data["prompts"].([]interface{}); ok {
			for _, tr := range templatesRaw {
				tmplData, _ := jsonMarshal(tr)
				var tmpl models.PromptTemplate
				if jsonUnmarshal(tmplData, &tmpl) == nil {
					if existing, _ := repo.GetTemplateByCode(tmpl.Code); existing != nil && existing.ID > 0 {
						tmpl.ID = existing.ID
						if err := repo.UpdateTemplate(&tmpl); err != nil {
							return err
						}
					} else {
						tmpl.ID = 0
						if err := repo.CreateTemplate(&tmpl); err != nil {
							return err
						}
					}
					imported["prompts"]++
				}
			}
		}

		return nil
	})

	if err != nil {
		respondInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Import successful", "imported": imported})
}
