import { Pressable, StyleSheet, Text } from 'react-native';

export default function Button({ label, onPress }) {
  return (
    <Pressable style={styles.button} onPress={onPress}>
      <Text style={styles.text}>
        {label}
      </Text>
    </Pressable>
  )
}

const styles = StyleSheet.create({
  button: {
    backgroundColor: '#4d676e',
    paddingVertical: 10,
    paddingHorizontal: 30,
    borderRadius: 9,
    color: 'white',
    fontWeight: 'bold',
  },
  text: {
    color: 'white',
  }
});
