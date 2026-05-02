import { ref } from 'vue'

export function useCrudDialog(initForm = {}) {
  const visible = ref(false)
  const loading = ref(false)
  const editing = ref(false)
  const form = ref({ ...initForm })

  function openCreate() {
    editing.value = false
    form.value = { ...initForm }
    visible.value = true
  }

  function openEdit(row) {
    editing.value = true
    form.value = { ...row }
    visible.value = true
  }

  function close() {
    visible.value = false
    loading.value = false
  }

  return { visible, loading, editing, form, openCreate, openEdit, close }
}
