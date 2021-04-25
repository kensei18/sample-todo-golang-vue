<template>
  <draggable
    group="task-list"
    style="min-height: 90vh"
    @add="$emit('update-status')"
  >
    <v-card
      v-for="task in taskList"
      :key="task.id"
      class="my-2"
      color="cyan lighten-5"
      @dragstart="updateSelectedTaskId(task.id)"
    >
      <v-card-title>
        <span>{{ task.name }}</span>
        <v-spacer />
        <v-btn fab text icon x-small @click="showEditDialog(task.id)">
          <v-icon>mdi-pencil</v-icon>
        </v-btn>
        <v-btn fab text icon x-small @click="removeTask(task.id)">
          <v-icon>mdi-delete</v-icon>
        </v-btn>
      </v-card-title>
      <v-card-text>{{ task.description }}</v-card-text>
    </v-card>
  </draggable>
</template>

<script>
import draggable from "vuedraggable";

export default {
  components: { draggable },

  props: {
    taskList: Array,
  },

  inject: ["updateSelectedTaskId", "showEditDialog", "removeTask"],
};
</script>
