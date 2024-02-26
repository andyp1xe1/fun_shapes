import { Pressable, StyleSheet, Text } from 'react-native';

export default function Button({ label, onPress }) {
  return (
    <Pressable style={styles.button} onPress={onPress}>
      <Text>
        {label}
      </Text>
    </Pressable>
  )
}

const styles = StyleSheet.create({
  button: {
    backgroundColor: '#D79921',
    paddingVertical: 10,
    paddingHorizontal: 30,
    borderRadius: 9,
  },
});
