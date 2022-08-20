package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/kakengloh/tsk/entity"
	"github.com/kakengloh/tsk/repository"
	"github.com/kakengloh/tsk/util/printer"
	"github.com/spf13/cobra"
)

func NewModCommand(tr repository.TaskRepository) *cobra.Command {
	setCmd := &cobra.Command{
		Use:   "mod",
		Short: "Modify an existing task",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pt := printer.New(cmd.OutOrStdout())

			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("task ID must be an integer: %w", err)
			}

			task, err := tr.GetTaskByID(id)
			if err != nil {
				return fmt.Errorf("task not found")
			}

			// Title
			title := task.Title
			t, err := cmd.Flags().GetString("title")
			if err != nil {
				return err
			}
			if t != "" {
				title = t
			}

			// Priority
			priority := task.Priority
			p, err := cmd.Flags().GetString("priority")
			if err != nil {
				return err
			}
			if p != "" {
				v, ok := entity.TaskPriorityFromString[p]
				if !ok {
					return fmt.Errorf("invalid priority: %s, valid values are [low, medium, high]", p)
				}
				priority = v
			}

			// Status
			status := task.Status
			s, err := cmd.Flags().GetString("status")
			if err != nil {
				return err
			}
			if s != "" {
				v, ok := entity.TaskStatusFromString[s]
				if !ok {
					return fmt.Errorf("invalid status: %s, valid values are [todo, doing, done]", s)
				}
				status = v
			}

			// Due
			d, err := cmd.Flags().GetString("due")
			if err != nil {
				return err
			}
			var due time.Time
			if d != "" {
				duration, err := time.ParseDuration(d)
				fmt.Println(err)
				if err == nil {
					due = time.Now().Add(duration)
				} else {
					val, err := time.ParseInLocation("2006-01-02 15:04", d, time.Local)
					if err != nil {
						return err
					}
					due = val
				}
			}

			task, err = tr.UpdateTask(id, title, priority, status, due)

			pt.PrintTask(task, "Task modified ✅")

			return err
		},
	}

	setCmd.PersistentFlags().StringP("title", "t", "", "Set title")
	setCmd.PersistentFlags().StringP("status", "s", "", "Set status (todo / doing / done)")
	setCmd.PersistentFlags().StringP("priority", "p", "", "Set priority (low / medium / high")
	setCmd.PersistentFlags().StringP("due", "d", "", "Set due")

	return setCmd
}
