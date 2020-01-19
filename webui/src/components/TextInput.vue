<template>
  <input
    type="text"
    :placeholder="placeholder"
    :value="value"
    @input="$emit('input', $event.target.value)"
    class="py-2 px-4 bg-white focus:outline-none rounded focus:shadow-outline border appearance-none leading-normal border-gray-200 bg-gray-100"
  />
</template>

<script>
export default {
  name: 'TextInput',
  props: ['value', 'placeholder'],
  model: {
    prop: 'value',
    event: 'input'
  }
}
</script>