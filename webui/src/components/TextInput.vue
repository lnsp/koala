<template>
  <input
    type="text"
    :placeholder="placeholder"
    :modelValue="modelValue"
    @input="$emit('update:modelValue', $event.target.value)"
    class="py-2 px-4 bg-white focus:outline-none rounded focus:shadow-outline border appearance-none leading-normal border-gray-200 bg-gray-100"
  />
</template>

<script>
export default {
  name: 'TextInput',
  props: ['modelValue', 'placeholder'],
  emits: ['update:modelValue']
}
</script>