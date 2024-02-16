import { StatusBar } from 'expo-status-bar';
import { StyleSheet, Text, View } from 'react-native';
import * as ImagePicker from 'expo-image-picker';

import Button from './Button';

import ImageViewer from './ImageViewer';
const PlaceholderImage = require('./assets/images/gruv.jpg')


//TODO: https://docs.expo.dev/tutorial/image-picker/#use-the-selected-image
export default function App() {
  const pickImageAsync = async () => {
    let result = await ImagePicker.launchImageLibraryAsync({
      allowsEditing: true,
      quality: 1,
    });

    if (!result.canceled) {
      console.log(result);
    } else {
      alert('You did not select any image.');
    }
  };

  return (
    <View style={styles.container}>
      <Text style={styles.heading}>
        Welcome to Fun Shapes!
      </Text>
      <ImageViewer placeholderImageSource={PlaceholderImage} />
      <Button label="choose your photo!" onPress={pickImageAsync} />
      <StatusBar style="auto" />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#282828',
    color: '#EBDBB2',
    alignItems: 'center',
    justifyContent: 'center',
    gap: 25,
  },
  heading: {
    color: '#EBDBB2',
    fontSize: 24,
    paddingVertical: 16,
    fontFamily: 'serif',
  },
});
