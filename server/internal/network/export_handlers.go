package network

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ddc-111/agentGame/server/internal/database/models"
	"github.com/ddc-111/agentGame/server/internal/database/repository"
)

// handleExport godoc
// @Summary      Export all data
// @Description  Export all game data as JSON
// @Tags         data
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /export [get]
func (s *Server) handleExport(c *gin.Context) {
	ctx := c.Request.Context()
	data := make(map[string]interface{})

	scenes, err := s.repo.GetScenes(ctx)
	if err != nil {
		respondInternalError(c, err)
		return
	}
	data["scenes"] = scenes

	npcs, err := s.repo.GetNPCs(ctx)
	if err != nil {
		respondInternalError(c, err)
		return
	}
	data["npcs"] = npcs

	agents, err := s.repo.GetAgents(ctx)
	if err != nil {
		respondInternalError(c, err)
		return
	}
	data["agents"] = agents

	shops, err := s.repo.GetShops(ctx)
	if err != nil {
		respondInternalError(c, err)
		return
	}
	data["shops"] = shops

	items, err := s.repo.GetItems(ctx)
	if err != nil {
		respondInternalError(c, err)
		return
	}
	data["items"] = items

	tasks, err := s.repo.GetTasks(ctx)
	if err != nil {
		respondInternalError(c, err)
		return
	}
	data["tasks"] = tasks

	flows, err := s.repo.GetFlows(ctx)
	if err != nil {
		respondInternalError(c, err)
		return
	}
	data["flows"] = flows

	templates, err := s.repo.GetTemplates(ctx)
	if err != nil {
		respondInternalError(c, err)
		return
	}
	data["prompts"] = templates

	c.JSON(http.StatusOK, gin.H{"data": data})
}

// handleImport godoc
// @Summary      Import data
// @Description  Import game data from JSON
// @Tags         data
// @Accept       json
// @Produce      json
// @Param        data  body  map[string]interface{}  true  "Import data"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /import [post]
func (s *Server) handleImport(c *gin.Context) {
	ctx := c.Request.Context()
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		respondError(c, http.StatusBadRequest, BadRequest(err.Error()))
		return
	}

	if inner, ok := data["data"].(map[string]interface{}); ok {
		data = inner
	}

	imported := make(map[string]int)

	err := s.repo.Transaction(ctx, func(repo *repository.Repository) error {
		if scenesRaw, ok := data["scenes"].([]interface{}); ok {
			for _, sr := range scenesRaw {
				sceneData, _ := jsonMarshal(sr)
				var scene models.Scene
				if jsonUnmarshal(sceneData, &scene) == nil {
					if existing, _ := repo.GetSceneByCode(ctx, scene.Code); existing != nil && existing.ID > 0 {
						scene.ID = existing.ID
						if err := repo.UpdateScene(ctx, &scene); err != nil {
							return err
						}
					} else {
						scene.ID = 0
						if err := repo.CreateScene(ctx, &scene); err != nil {
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
					if existing, _ := repo.GetNPCByCode(ctx, npc.Code); existing != nil && existing.ID > 0 {
						npc.ID = existing.ID
						if err := repo.UpdateNPC(ctx, &npc); err != nil {
							return err
						}
					} else {
						npc.ID = 0
						if err := repo.CreateNPC(ctx, &npc); err != nil {
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
					if existing, _ := repo.GetAgentByCode(ctx, agent.Code); existing != nil && existing.ID > 0 {
						agent.ID = existing.ID
						if err := repo.UpdateAgent(ctx, &agent); err != nil {
							return err
						}
					} else {
						agent.ID = 0
						if err := repo.CreateAgent(ctx, &agent); err != nil {
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
					if existing, _ := repo.GetShopByCode(ctx, shop.Code); existing != nil && existing.ID > 0 {
						shop.ID = existing.ID
						if err := repo.UpdateShop(ctx, &shop); err != nil {
							return err
						}
					} else {
						shop.ID = 0
						if err := repo.CreateShop(ctx, &shop); err != nil {
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
					if existing, _ := repo.GetItemByCode(ctx, item.Code); existing != nil && existing.ID > 0 {
						item.ID = existing.ID
						if err := repo.UpdateItem(ctx, &item); err != nil {
							return err
						}
					} else {
						item.ID = 0
						if err := repo.CreateItem(ctx, &item); err != nil {
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
					if existing, _ := repo.GetTaskByCode(ctx, task.Code); existing != nil && existing.ID > 0 {
						task.ID = existing.ID
						if err := repo.UpdateTask(ctx, &task); err != nil {
							return err
						}
					} else {
						task.ID = 0
						if err := repo.CreateTask(ctx, &task); err != nil {
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
					if existing, _ := repo.GetFlowByCode(ctx, flow.Code); existing != nil && existing.ID > 0 {
						flow.ID = existing.ID
						if err := repo.UpdateFlow(ctx, &flow); err != nil {
							return err
						}
					} else {
						flow.ID = 0
						if err := repo.CreateFlow(ctx, &flow); err != nil {
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
					if existing, _ := repo.GetTemplateByCode(ctx, tmpl.Code); existing != nil && existing.ID > 0 {
						tmpl.ID = existing.ID
						if err := repo.UpdateTemplate(ctx, &tmpl); err != nil {
							return err
						}
					} else {
						tmpl.ID = 0
						if err := repo.CreateTemplate(ctx, &tmpl); err != nil {
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
